# musiman
A tool to manage your music files

## Get musiman
TBD

## Usage
TBD

## Features
- [x] recursively scan a directory for music files 
- [ ] add minimum size filter to ignore tiny audio files from i.e. game sound effects
- [ ] add path ignore patterns (exact relative paths for now)
- [ ] store music files in sqlite with calculated content hash (NOT acustid, just a hash)
- [ ] decode audio files (mp3 only for now) to get raw audio
- [ ] integrate [gochroma](https://github.com/go-fingerprint/gochroma) to get acustid (audio fingerprint)
- [ ] store acustids for files in sqlite
- [ ] lookup [musicbrainz](https://musicbrainz.org/) data by [acustid](https://acoustid.org/)
- [ ] store [musicbrainz](https://musicbrainz.org/) data for files in sqlite
- [ ] read/write metadata from and to files
- [ ] deduplicate audio files based on hash and acustid (always keeps the best quality version)
- [ ] convert audio file formats
- [ ] create a central media library
  > A central folder for all (deduped) music, optionally converted to a unified file format, with a filesystem hierarchy like:
  ```
   <Artist>
     <Album>
       <Track>
       <Track>
       ...
     <Album>
     ...
   <Artist>
   ...
  ```

## Build from source
TBD