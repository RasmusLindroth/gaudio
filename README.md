# gaudio

A quick and dirty hack with pulseaudio to change the current active soundcard
with [gtoggle](./gtoggle) and change the volume with [gvol](./gvol).

I use it with i3 so I can change all this from the keyboard. My config: 

```
# Change volume
bindsym XF86AudioRaiseVolume exec --no-startup-id gvol up 5 
bindsym XF86AudioLowerVolume exec --no-startup-id gvol down 5

# Toogle between audio outputs
bindcode 202 exec --no-startup-id gtoggle alsa_output.pci-0000_00_1b.0.analog-stereo alsa_output.usb-Kingston_HyperX_7.1_Audio_00000000-00.analog-stereo
```
