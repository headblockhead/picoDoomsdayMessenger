# PCB Files
The files in this directory can be used to manufacture a copy of the PCB used in this project. They are here as a backup in case EASYEDA is ever taken down. The files are also here in case you want to make a copy of the PCB for yourself.
You can find the files online at this [open source hardware lab page](https://oshwlab.com/headblockhead/picodoomsdaycommunicator).

## Latest PCB Version: V1.2

## PCB Design
The PCB was designed using [EasyEDA](https://easyeda.com/). The design files can be found in the [EasyEDA project](https://oshwlab.com/headblockhead/picodoomsdaycommunicator). The design files are also in this directory.

## V1 vs V1.1
The difference between V1 and V1.1 is that the logo and text on the PCB have been moved to the silkscreen layer. In V1, the logo and text were on the soldermask layer. This meant that the logo and text were hard to see in low light conditions. The logo and text are now on the silkscreen layer which means that they are easier to see in low light conditions.

## V1.1 vs V1.2
The difference between V1.1 and V1.2 is that in the previous two designs, there was no wire between switch7 and the appropriate diode. This was due to an error in the naming of the net. The switch's net was named "SWITCH7" and the diode's net was named "SWITCH_7". This caused the router to not connect them.
