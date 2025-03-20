package erro

import "errors"

var (
	ErrorNotPost          = errors.New("Request is not POST")
	ErrorNotGet           = errors.New("Request is not GET")
	ErrorNullEmail        = errors.New("Email is Null")
	ErrorNotEmail         = errors.New("Email is not validate")
	ErrorNullName         = errors.New("Name is Null")
	ErrorNotReadAll       = errors.New("Data is not readAll")
	ErrorUnmarshal        = errors.New("Unmarshal's error")
	ErrorMarshal          = errors.New("Marshal's error")
	ErrorValidator        = errors.New("Validator's error")
	ErrorResponse         = errors.New("Response's error")
	ErrorGetEnv           = errors.New("Environment's error")
	ErrorDBConnect        = errors.New("Connect to DB failed")
	ErrorDBPing           = errors.New("DB-Ping failed")
	ErrorDBLS             = errors.New("ListenAndServe error")
	ErrorServerShutdown   = errors.New("Server Shutdown error")
	ErrorDBAdd            = errors.New("DB's Add error")
	ErrorDBGet            = errors.New("DB's Get error")
	ErrorHash             = errors.New("Hash error!")
	ErrorUniqueEmail      = errors.New("This email has already been registered")
	ErrorEmailNotRegister = errors.New("This email is not registered")
	ErrorInvalidPerson    = errors.New("Invalid password!")
)
