package config

import "github.com/gorilla/sessions"

const SESSIONS_ID = "go_reg_sessions"

var Store = sessions.NewCookieStore([]byte("jkasdfsad1243bkjasdf"))