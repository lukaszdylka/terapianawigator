# Kalendarz Specjalisty

Szyfrowany kalendarz pracy psychologa, pedagoga i innych specjalistów szkolnych.

## Wersja docelowa: PWA

Aplikacja działa jako PWA pod adresem:

`https://terapianawigator.pl/kalendarz-specjalisty/`

Najlepiej otworzyć ten adres w Brave i zainstalować aplikację z poziomu Brave. Po instalacji uruchamia się w osobnym oknie bez kart i paska adresu.

Nie jest wymagany instalator EXE ani Microsoft Edge.

## Aktualizacje

Kod PWA jest pobierany z serwera przy uruchomieniu. Service Worker przechowuje ostatnią działającą wersję do pracy offline.

Aktualizacja nie podmienia plików na komputerze i nie uruchamia skryptów CMD ani PowerShell. Dzięki temu nie ma problemu z aktualizacją EXE, polskimi znakami w ścieżkach ani blokadą pliku przez Windows.

## Dane

Dane uczniów i spotkań:
- pozostają lokalnie w profilu przeglądarki Brave,
- są szyfrowane hasłem,
- nie są zapisywane w repozytorium ani bazie online,
- można eksportować i importować przez zaszyfrowane kopie `.ksbackup`.

Przy przejściu ze starego EXE do PWA należy wyeksportować kopię w starej aplikacji i zaimportować ją do PWA, ponieważ `127.0.0.1` i `terapianawigator.pl` mają osobne magazyny przeglądarki.

## Offline

Po pierwszym poprawnym uruchomieniu online aplikacja działa także bez Internetu dzięki Service Workerowi.

## Windows Hello

Jeżeli Brave i komputer obsługują wymagane mechanizmy WebAuthn i Windows Hello, w ustawieniach można włączyć odblokowanie biometrią albo PIN-em Windows. Hasło szyfrowania pozostaje metodą awaryjną.

## Stara wersja EXE

Pliki Go pozostają w repozytorium jako wersja archiwalna. Automatyczne publikowanie kolejnych EXE zostało wyłączone.
