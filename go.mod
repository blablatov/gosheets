module github.com/blablatov/gosheets

go 1.16

replace github.com/blablatov/gosheets/factpost => ./factpost

replace github.com/blablatov/gosheets/getmo => ./getmo

replace github.com/blablatov/gosheets/planpost => ./planpost

replace github.com/blablatov/gosheets/weipost => ./weipost

require (
	golang.org/x/oauth2 v0.2.0
	google.golang.org/api v0.103.0
)
