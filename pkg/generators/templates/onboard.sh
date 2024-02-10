#!/bin/sh
# shellcheck disable=SC1083
#
# Perform early onboarding of BIG-IP for Google Cloud

# Log an info message
info()
{
    echo "$0: INFO: $*" >&2
}

# Log an error message and exit
error()
{
    echo "$0: ERROR: $*" >&2
    exit 1
}

# Return an bearer token for the VM service account to use with GCP APIs
auth_token()
{
    attempt=0
    while [ "${attempt}" -lt 10 ]; do
        auth_token="$(curl -sf --retry 20 -H "Metdata-Flavor: Google" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token | jq --raw-output '.access_token')"
        retval=$?
        if [ "${retval}" -eq 0 ]; then
            echo "${auth_token}"
            break
        fi
        info "auth_token: ${attempt}: Curl failed with exit code $?; sleeping before retry"
        sleep 15
        attempt=$((attempt+1))
    done
    [ "${attempt}" -ge 10 ] && \
        info "auth_token: ${attempt}: Failed to get an auth token from metadata server"
    # shellcheck disable=SC2086
    return ${retval}
}

# Download the remote resource to provided path
# $1 = URL
# $2 = output file
# $3+ are additional curl arguments
retry_download()
{
    url="$1"
    out="$2"
    shift
    shift
    attempt=0
    while [ "${attempt}" -lt 10 ]; do
        info "retry_download: ${attempt}: Downloading ${url} to ${out}"
        curl -sfL --retry 20 {{ with .ProxyURL }}-x "{{ . }}"{{ end -}} -o "${out}" "$@" "${url}"
        retval=$?
        [ "${retval}" -eq 0 ] && break
        info "retry_download: ${attempt}: Failed to download ${url}: exit code: ${retval}; sleeping before retrying"
        sleep 15
        attempt=$((attempt+1))
    done
    [ "${attempt}" -ge 10 ] && \
        info "retry_download: Failed to download from ${url}; giving up"
    # shellcheck disable=SC2086
    return ${retval}
}

# Recognise GCS storage API requests and handle authentication as necessary.
# $1 = URL of remote file
# $2 = output path
download()
{
    mkdir -p "$(dirname "$2")" || \
        error "Error creating directory for $2; exit code $?"
    case "$1" in
        gs://*)
            gs_uri="$(printf '%s' "${1}" | jq --slurp --raw-input --raw-output 'split("/")[2:]|["", .[0], "/o/", (.[1:]|join("/")|@uri), "?alt=media"]|join("")')" || \
                error "Error creating JSON API URL from ${1}; exit code $?"
            auth_token="$(auth_token)" || error "Unable to get auth token"
            retry_download "${gs_uri}" "$2" -H "Authorization: Bearer ${auth_token}"
            ;;
        https://storage.googleapis.com/*)
            auth_token="$(auth_token)" || error "Unable to get auth token"
            retry_download "$1" "$2" -H "Authorization: Bearer ${auth_token}"
            ;;
        ftp://*|http://*|https://*)
            retry_download "$1" "$2"
            ;;
        /*)
            cp "$1" "$2"
            ;;
        *)
            info "Unrecognised remote scheme for $1"
            false
            ;;
    esac
    return $?
}

mkdir -p /var/config/rest/downloads

info "Starting to onboard; waiting for BIG-IP to be ready"

# shellcheck disable=SC1091
. /usr/lib/bigstart/bigip-ready-functions
wait_bigip_ready
nic_count="$(curl -sf --retry 20 -H "Metadata-Flavor: Google" "http://169.254.169.254/computeMetadata/v1/instance/network-interfaces/?recursive=true" | jq --raw-output '.|length')"
if [ "${nic_count}" -gt 1 ] && [ ! -f /config/cloud/.mgmtInterface ]; then
    info "Getting management interface configuration from metadata"
    target_address="$(curl -sf --retry 20 -H "Metadata-Flavor: Google" http://169.254.169.254/computeMetadata/v1/instance/network-interfaces/1/ip)"
    target_netmask="$(curl -sf --retry 20 -H "Metadata-Flavor: Google" http://169.254.169.254/computeMetadata/v1/instance/network-interfaces/1/subnetmask)"
    target_gateway="$(curl -sf --retry 20 -H "Metadata-Flavor: Google" http://169.254.169.254/computeMetadata/v1/instance/network-interfaces/1/gateway)"
    target_mtu="$(curl -sf --retry 20 -H "Metadata-Flavor: Google" http://169.254.169.254/computeMetadata/v1/instance/network-interfaces/1/mtu)"
    target_network="$(ipcalc -n "${target_address}" "${target_netmask}" | cut -d= -f2)"
    # NOTE: this configuration is based on f5devcentral/terraform-gcp bigip-module boarding script, unless called out
    info "Resetting management interface."
    tmsh modify sys global-settings gui-setup disabled
    tmsh modify sys global-settings mgmt-dhcp disabled
    tmsh delete sys management-route all
    tmsh delete sys management-ip all
    info "Configuring management interface"
    tmsh create sys management-ip "${target_address}/32"
    tmsh create sys management-route mgmt_gw network "${target_gateway}/32" type interface mtu "${target_mtu}"
    tmsh create sys management-route mgmt_net network "${target_network}/${target_netmask}" gateway "${target_gateway}" mtu "${target_mtu}"
    tmsh create sys management-route default gateway "${target_gateway}" mtu "${target_mtu}"
    tmsh modify sys global-settings remote-host add { metadata.google.internal { hostname metadata.google.internal addr 169.254.169.254 } }
    tmsh modify sys management-dhcp sys-mgmt-dhcp-config request-options delete { ntp-servers }
    # MEmes - make sure the GCP metadata server is used for DNS, at least until user explicitly overrides in DO
    tmsh modify sys dns name-servers add { 169.254.169.254 }
    tmsh save /sys config
    touch /config/cloud/.mgmtInterface
    info "Setup of management interface is complete."
fi

# Is a management NIC swap necessary?
current_mgmt_nic="$(tmsh list sys db provision.managementeth value 2>/dev/null | awk -F\" 'NR==2 {print $2}')"
if [ "${nic_count}" -gt 1 ] && [ "${current_mgmt_nic}" != "eth1" ]; then
    info "Management NIC swap is necessary; updating database."
    bigstart stop tmm
    tmsh modify sys db provision.managementeth value eth1
    tmsh modify sys db provision.1nicautoconfig value disable
    tmsh save /sys config
    [ -e "/etc/ts/common/image.cfg" ] && \
        sed -i "s/iface=eth0/iface=eth1/g" /etc/ts/common/image.cfg
    info "Rebooting for management NIC swap."
    reboot
    exit 0
else
    info "Management NIC swap is not needed; continuing"
fi

{{- with .RuntimeInit }}
# Download and execute runtime-init
info "Downloading runtime-init installer from {{ .PackageURL }}"
download "{{ .PackageURL }}" "/var/config/rest/downloads/f5-bigip-runtime-init.gz.run" || \
    error "Failed to download {{ .PackageURL }}: exit code: $?"
echo "{{ .PackageSHA }} /var/config/rest/downloads/f5-bigip-runtime-init.gz.run" | \
    sha256sum --status --check || \
    error "Failed to verify integrity of {{ .PackageURL }}: exit code: $?"

if [ ! -x /usr/local/bin/f5-bigip-runtime-init ] && [ -f /var/config/rest/downloads/f5-bigip-runtime-init.gz.run ]; then
    info "Installing runtime-init package"
    bash /var/config/rest/downloads/f5-bigip-runtime-init.gz.run -- '--cloud gcp' || \
        error "Failed to install runtime-init: exit code $?"
fi

if [ -f /config/cloud/runtime-init-conf.yaml ]; then
    if [ -x /usr/local/bin/f5-bigip-runtime-init ]; then
        info "Executing runtime-init"
        /usr/local/bin/f5-bigip-runtime-init --config-file /config/cloud/runtime-init-conf.yaml
        retval=$?
        if [ "${retval}" -ne 0 ]; then
            cat /var/log/cloud/bigIpRuntimeInit.log
            error "Failed to execute runtime-init: exit code: ${retval}"
        fi
    else
        error "Runtime-init is not installed; skipping"
    fi
else
    error "Runtime-init configuration file was not found at /config/cloud/runtime-init-conf.yaml"
fi
{{- end }}

# Disable the onboarding systemd unit and enable management route reset unit for
# future boots
info "Onboarding complete; disabling f5-gce-onboard.service unit"
systemctl disable f5-gce-onboard.service
info "Onboarding complete; enable f5-gce-management-route.service unit"
systemctl enable f5-gce-management-route.service

info "Onboarding complete."
