# hue-state-cyler
This is a tiny microservice which runs on our resident k3s cluster.
It is used to power cycle the TTN (The Things Network) base station
which is hooked up to a Hue smart plug. 

# rationale
The TTN base station or the TTN network itself have become very 
flaky recently. This requires a power cycle every 2-3 days of the base station.
Therefore we hooked up the TTN base to a Hue smart plug which can be controlled
by this microservice. A Grafana alert triggers the hue-state-cycler in case packets
 are not received for 5 consequitve minutes. A device can only be cycled every 
5 minutes. This prevents continuous cycling in case there are network or other external
issues.

# license
License is MIT
