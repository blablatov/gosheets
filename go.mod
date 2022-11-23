module github.com/blablatov/gosheets

go 1.16

replace github.com/blablatov/gosheets/factpost => ./factpost

replace github.com/blablatov/gosheets/getmo => ./getmo

replace github.com/blablatov/gosheets/planpost => ./planpost

replace github.com/blablatov/gosheets/weipost => ./weipost

require (
	github.com/pion/stun v0.3.5 // indirect
	golang.org/x/oauth2 v0.2.0
	google.golang.org/api v0.103.0
	google.golang.org/genproto v0.0.0-20221027153422-115e99e71e1c // indirect
)
