# SPDX-FileCopyrightText: 2021 ladyada for Adafruit Industries
# SPDX-License-Identifier: MIT

# Example to send a packet periodically between addressed nodes
# Author: Jerry Needell
# Modified to work directly on the Challenger RP2040 LoRa boards.
#
import time
import board
import busio
import digitalio
import adafruit_rfm9x

import displayio
import terminalio
import label
import adafruit_displayio_ssd1306

displayio.release_displays()

# Use for I2C
i2c = busio.I2C(scl=board.GP1, sda=board.GP0)
display_bus = displayio.I2CDisplay(i2c, device_address=0x3C)

WIDTH = 128
HEIGHT = 64
BORDER = 5

display = adafruit_displayio_ssd1306.SSD1306(display_bus, width=WIDTH, height=HEIGHT)

# set the time interval (seconds) for sending packets
transmit_interval = 1

# Define radio parameters.
RADIO_FREQ_MHZ = 868.0  # Frequency of the radio in Mhz. Must match your
# module! Can be a value like 915.0, 433.0, etc.

# Define pins connected to the chip.
CS = digitalio.DigitalInOut(board.RFM95W_CS)
RESET = digitalio.DigitalInOut(board.RFM95W_RST)

led = digitalio.DigitalInOut(board.GP24)
led.direction = digitalio.Direction.OUTPUT

# Initialize SPI bus.
rfm95x_spi = busio.SPI(board.RFM95W_SCK, MOSI=board.RFM95W_SDO, MISO=board.RFM95W_SDI)
# Initialze RFM radio
rfm9x = adafruit_rfm9x.RFM9x(rfm95x_spi, CS, RESET, RADIO_FREQ_MHZ)

# set node addresses
# When trying this out with two boards, one board should have the reverse settings
rfm9x.node = 1
rfm9x.destination = 2
# initialize counter
counter = 0
# send a broadcast message from my_node with ID = counter
rfm9x.send(
    bytes("Startup message {} from node {}".format(counter, rfm9x.node), "UTF-8")
)

# Wait to receive packets.
print("Waiting for packets...")
now = time.monotonic()

led.value = True

while True:
    # Make the display context
    splash = displayio.Group()
    display.show(splash)

    color_bitmap = displayio.Bitmap(WIDTH, HEIGHT, 1)
    color_palette = displayio.Palette(1)
    color_palette[0] = 0x000000

    bg_sprite = displayio.TileGrid(color_bitmap, pixel_shader=color_palette, x=0, y=0)
    splash.append(bg_sprite)
    
    led.value = False
    
    # Look for a new packet: only accept if addresses to my_node
    packet = rfm9x.receive(with_header=True)
    # If no packet was received during the timeout then None is returned.
    if packet is not None:
        # Received a packet!
        # Print out the raw bytes of the packet:
        print("Received (raw header):", [hex(x) for x in packet[0:4]])
        print("Received (raw payload): {0}".format(packet[4:]))
        text3 = "{0}".format(bytes(packet[4:]))
        text3_area = label.Label(
            terminalio.FONT, text=text3, color=0xFFFFFF, x=0, y=10
        )
        splash.append(text3_area)
        print("Received RSSI: {0}".format(rfm9x.last_rssi))
        text2 = "RSSI: {0}".format(rfm9x.last_rssi)
        text2_area = label.Label(
            terminalio.FONT, text=text2, color=0xFFFFFF, x=0, y=25
        )
        splash.append(text2_area)
        display.refresh()
        led.value = True
        time.sleep(0.5)
    if time.monotonic() - now > transmit_interval:
        now = time.monotonic()
        counter = counter + 1
        # send a  mesage to destination_node from my_node
        rfm9x.send(
            bytes(
                "msg num {} node {}".format(counter, rfm9x.node), "UTF-8"
            ),
            keep_listening=True,
        )
        button_pressed = None


