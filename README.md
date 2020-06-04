# SML smart meter reader & parser in go(lang)
!!This is a PoC!! and is currently working only with the following smart meter -> https://www.emh-metering.de/produkte/smart-meter/ehz-k

Nevertheless it should be easy to extend/change the implementation to work with other smart meters. I didn't find a  generally applicable documentation of the SML, so I had to test and analyse the return value of my smart meter.
There are some good sites (only in german - sorry) where you get some explanations about the SML (also for other smart meters). If you have questions or improvement ideas feel free to ask or create a pull request.
This were some of my first go steps, architecture and code is far from perfect :) .

## Link collection
Java-Library zum analysieren der Daten:
https://mvnrepository.com/artifact/org.openmuc/jsml/1.1.2
https://www.openmuc.org/sml/
https://github.com/jblu48317/SMLToJSON
https://linuxize.com/post/install-java-on-raspberry-pi/

Über SML:
https://de.wikipedia.org/wiki/Smart_Message_Language


Generelle Seiten:
http://www.stefan-weigert.de/php_loader/sml.php
https://www.rudiswiki.de/wiki9/VolkszaehlerEMHeHZ
https://www.msxfaq.de/sonst/bastelbude/smartmeter_d0_sml.htm
https://www.msxfaq.de/sonst/bastelbude/smartmeter_d0_sml_protokoll.htm

http://blog.bubux.de/raspberry-pi-ehz-auslesen/

https://wiki.volkszaehler.org/hardware/channels/meters/power/edl-ehz/emh-ehz-h1
https://wiki.volkszaehler.org/hardware/channels/meters/power/edl-ehz/edl21-ehz

https://www.emh-metering.de/produkte/smart-meter/ehz-k
https://www.emh-metering.de/images/Produkt-Dokumentation/eHZ-K-BIA-D-1-20.pdf

https://www.bsi.bund.de/SharedDocs/Downloads/DE/BSI/Publikationen/TechnischeRichtlinien/TR03109/TR-03109-1_Anlage_Feinspezifikation_Drahtgebundene_LMN-Schnittstelle_Teilb.pdf?__blob=publicationFile
http://itrona.ch/stuff/F2-2_PJM_5_Beschreibung%20SML%20Datenprotokoll%20V1.0_28.02.2011.pdf
https://wiki.volkszaehler.org/hardware/channels/meters/power/edl-ehz/emh-ehz-k#beispielkonfiguation
https://wiki.volkszaehler.org/software/sml#beispiel_3emh_ehz_fw8e2a50bak2


# BUILD FOR RASPBERRY PI (RASPBERIAN)
env GOOS=linux GOARCH=arm GOARM=5 go build