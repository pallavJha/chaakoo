package chaakoo

import (
  "os"
  "path/filepath"
)


type SavedSession struct {
	SessionName string
	configFile  string
}

type SavedContent struct {
  Sessions []SavedSession
}



func loadSessions() []SavedSession {
  sessionPath := getSavedSessionFilePath()
  contents, err := os.ReadFile(sessionPath)
  if err != nil {
    if os.IsNotExist(err) {
      f,err := os.Create(sessionPath)
      if err != nil {
        log.Fatal().Err(err).Msg("cannot create the file for saving the sessions")
      } else {
        defer func() {
          err := f.Close()
          if err != nil {
            log.Debug().Err(err).Msg("cannot close the session file")
          }
        }
      }
    } else {
      log.Fatal().Err(err).Msg("cannot read the sessions file from the path")
    }
  }
  return []
}


func getSavedSessionFilePath() string {
  dirname, err := os.UserHomeDir()
  if err != nil {
    log.Fatal().Err(err).Msg("cannot access the user's home directory")
  }

  return filepath.Join(dirname, ".config", ".chaakoo")
}
