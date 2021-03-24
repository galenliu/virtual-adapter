from gateway_addon import Device
from pkg.pybuspro.devices.light import Light
import threading
import time

from pkg.buspro_property import BusproProperty

_POLL_INTERVAL = 5


class BusproDevice(Device):

    def __init__(self, adapter, device_address, channel_number, title=""):
        """
        Initialize the object.
        adapter -- the Adapter managing this device
        _id -- ID of this device
        hs100_dev -- the pyHS100 device object to initialize from
        index -- index inside parent device
        """

        self._device_address = device_address
        device_id = device_address + channel_number
        Device.__init__(self, adapter, device_id)
        if title == "":
            self.title = self._type[0] + device_id
        self._channel_number = channel_number


class BusproLight(BusproDevice):
    """TP-Link smart plug type."""

    def __init__(self, adapter, sub_id, device_id, channel_number, is_dimmable=False, is_color=False):
        """
        Initialize the object.
        adapter -- the Adapter managing this device
        _id -- ID of this device
        hs100_dev -- the pyHS100 device object to initialize from
        index -- index inside parent device
        """
        BusproDevice.__init__(self,
                              adapter,
                              device_id,
                              channel_number,
                              )
        self._type.append('Light')

        self.is_dimmable = is_dimmable
        self.is_color = is_color
        self.light_device = Light(adapter.gateway, device_address=(sub_id, device_id), channel_number=channel_number)
        self.light_device.register_device_updated_cb(self.poll)
        if self.is_dimmable:
            self.properties['level'] = BusproProperty(
                self,
                'level',
                {
                    '@type': 'BrightnessProperty',
                    'title': 'Brightness',
                    'type': 'integer',
                    'unit': 'percent',
                    'minimum': 0,
                    'maximum': 100,
                },
                100
            )

        if self.is_color:
            self._type.append('ColorControl')
            self.properties['color'] = BusproProperty(
                self,
                'color',
                {
                    '@type': 'ColorProperty',
                    'title': 'Color',
                    'type': 'string',
                },
                ""

            )

        self.properties['on'] = BusproProperty(
            self,
            'on',
            {
                '@type': 'OnOffProperty',
                'title': 'On/Off',
                'type': 'boolean',
            },
            ""
        )

    def poll(self, data):
        print(data)
        return
