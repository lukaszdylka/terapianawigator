# Kalendarz Specjalisty

Offline'owy kalendarz pracy psychologa, pedagoga i innych specjalistów szkolnych.

## Dane

Aplikacja nie korzysta z bazy danych online. Dane uczniów i spotkań są przechowywane lokalnie na komputerze użytkownika i szyfrowane hasłem.

## Budowanie ręczne

W katalogu `kalendarz-specjalisty`:

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-H=windowsgui -s -w" -o KalendarzSpecjalisty.exe .
```

## Aktualizacje

Aplikacja sprawdza manifest:

`https://raw.githubusercontent.com/lukaszdylka/terapianawigator/main/kalendarz-specjalisty/update.json`

Plik `update.json` jest aktualizowany automatycznie przez GitHub Actions po zmianach w kodzie aplikacji.

Workflow:
- buduje nowe EXE,
- oblicza SHA-256,
- publikuje Release na GitHubie,
- aktualizuje `update.json`,
- dzięki temu zainstalowane aplikacje mogą pobrać nową wersję bez ręcznego podmieniania plików.

## Kopie zapasowe

Automatyczne zaszyfrowane kopie są zapisywane lokalnie w:

`%LOCALAPPDATA%\KalendarzSpecjalisty\Backups`

Pliki `.ksbackup` pozostają zaszyfrowane tym samym hasłem co dane aplikacji.
