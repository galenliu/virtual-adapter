"""TP-Link adapter for WebThings Gateway."""

from gateway_addon import Adapter, Database

from pkg.buspro_device import BusproLight
from pkg.pybuspro.buspro import Buspro
from pyHS100 import Discover, SmartBulb, SmartPlug, SmartStrip

_TIMEOUT = 3


class BusproAdapter(Adapter):
    """Adapter for TP-Link smart home devices."""

    def __init__(self, verbose=False):
        """
        Initialize the object.
        verbose -- whether or not to enable verbose logging
        """
        self.name = self.__class__.__name__
        Adapter.__init__(self,
                         'buspro-adapter',
                         'buspro-adapter',
                         verbose=verbose)

        self.pairing = False
        self.start_pairing(_TIMEOUT)

    def _add_from_config(self):
        """Attempt to add all configured devices."""

        database = Database('buspro-adapter')
        if not database.open():
            return

        config = database.load_config()
        database.close()

        if not config or 'addresses' not in config:
            return

        ip = config["ip"]
        port = config["ip"]
        gateway_address_send_receive = ((ip, port), ('', port))
        self.gateway = Buspro(gateway_address_send_receive)

        for address in config['addresses']:

            try:
                if address["type"] == "Colour Light":
                    dev = BusproLight(self, address["subnet_id"], address["ip"], address["port"], is_color=True)
                    self.handle_device_added(dev)

                if address["type"] == "Dimmer":
                    dev = BusproLight(self, address["subnet_id"], address["ip"], address["port"], is_dimmable=True)
                    self.handle_device_added(dev)

                if address["type"] == "Light":
                    dev = BusproLight(self, address["subnet_id"], address["ip"], address["port"])
                    self.handle_device_added(dev)

            except (OSError, UnboundLocalError) as e:
                print('Failed to connect to {}: {}'.format(address, e))
                continue

    def start_pairing(self, timeout):
        """
        Start the pairing process.
        timeout -- Timeout in seconds at which to quit pairing
        """
        if self.pairing:
            return
        self.pairing = True
        self._add_from_config()
        self.pairing = False

    def cancel_pairing(self):
        """Cancel the pairing process."""
        self.pairing = False
