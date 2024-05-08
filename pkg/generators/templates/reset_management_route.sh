#!/bin/sh
#
# This script is intended to be executed on every boot to ensure the management
# default gateway is set correctly. See also `f5-gce-management-route.service`
# unit on instance.
#
# https://support.f5.com/csp/article/K85730674

info()
{
    echo "$0: INFO: $*" >&2
}

info "Management route reset handler: starting, waiting for BIG-IP to be ready"

# shellcheck disable=SC1091
. /usr/lib/bigstart/bigip-ready-functions
wait_bigip_ready

nic_count="$(curl -sf --retry 20 -H "Metadata-Flavor: Google" "http://169.254.169.254/computeMetadata/v1/instance/network-interfaces/?recursive=true" | jq --raw-output '.|length')"
if [ "${nic_count}" -gt 1 ]; then
    target_gateway="$(curl -sf --retry 20 -H "Metadata-Flavor: Google" http://169.254.169.254/computeMetadata/v1/instance/network-interfaces/1/gateway)"
    target_mtu="$(curl -sf --retry 20 -H "Metadata-Flavor: Google" http://169.254.169.254/computeMetadata/v1/instance/network-interfaces/1/mtu)"
    current_gw="$(tmsh list sys management-route default gateway | awk 'NR==2 { print $2 }')"
    while [ "${current_gw}" != "${target_gateway}" ]; do
        info "Management route reset handler: setting default gateway to ${target_gateway}; was ${current_gw}."
        tmsh delete sys management-route default
        tmsh create sys management-route default gateway "${target_gateway}" mtu "${target_mtu}"
        tmsh save /sys config
        current_gw="$(tmsh list sys management-route default gateway | awk 'NR==2 { print $2 }')"
    done
    info "Management route reset handler: complete."
else
    info "Managment route reset handler: nothing to do"
fi
