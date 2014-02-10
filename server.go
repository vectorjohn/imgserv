package main

import (
    "io"
    "io/ioutil"
	"net/http"
    "bytes"
	"os"
	"fmt"
    "log"
    //"path/filepath"
    "strings"
    //"image/draw"
    "image"
    "image/jpeg"
    _ "image/png"
    "code.google.com/p/graphics-go/graphics"
)

func main() {
	port := "4001"
	if len( os.Args ) > 1 {
		port = os.Args[1]
	}

	http.HandleFunc( "/", func( resp http.ResponseWriter, req *http.Request ) {
        //vals := req.URL.Path
        io.WriteString( resp, "Hello world!\n" + req.URL.Path )
    })

    relHandler( "/thumb/", serveThumb )

	listenOn := "localhost:" + port
	fmt.Println( "listening on " + listenOn )
	
	http.ListenAndServe(listenOn, nil )
}


var icache map[string]*bytes.Reader

func serveThumb( resp http.ResponseWriter, req *http.Request, path string ) {

    var outerr error

    log.Println( "Requested image", path )

    defer func() {
        //if err := recover(); err != nil {
        if recover() != nil {
            if outerr == nil {
                http.Error( resp, "Internal Server Error", 500 )
            } else {
                http.Error( resp, outerr.Error(), 500 )
            }
        }
    }()

    if icache == nil {
        icache = make( map[string]*bytes.Reader, 100 )
    }

    cached, ok := icache[ path ]

    if !ok {
        log.Println( "Generating thumbnail" )
        fd, err := os.Open( path )
        defer fd.Close()

        if err != nil {
            outerr = err;
            log.Panic( "ERROR: file open error ", err.Error() )
        }

        img, _, err := image.Decode( fd )
        if err != nil {
            outerr = err;
            log.Panic( "ERROR: Could not decode. ", err.Error() )
        }

        scaled := image.NewRGBA( image.Rect( 0, 0, 100, 100 ) )

        graphics.Thumbnail( scaled, img )

        var buf bytes.Buffer
        err = jpeg.Encode( &buf, scaled, &jpeg.Options{jpeg.DefaultQuality} )

        if  err != nil {
            outerr = err;
            log.Panic( "ERROR: could not encode to jpeg. ", err.Error() )
        }

        data, err := ioutil.ReadAll( &buf )
        if err != nil {
            outerr = err;
            log.Panic( "ERROR: ", err.Error() )
        }

        cached = bytes.NewReader( data )
        icache[ path ] = cached
    }

    cached.Seek( 0, 0 )
    io.Copy( resp, cached ) 
}

func relHandler( prefix string, handler func( http.ResponseWriter, *http.Request, string ) ) {
    // -1 for the trailing slash.  /foo/ should be 2
    prefLen := len( strings.Split( prefix, "/" ) ) - 1
    http.HandleFunc( prefix, relServer( handler, prefLen ) )
}

func relServer( handler func( http.ResponseWriter, *http.Request, string ), prefixLength int ) ( func( http.ResponseWriter, *http.Request ) ) {
    return func( resp http.ResponseWriter, req *http.Request ) {
        suffix := pathStripPrefix( req.URL.Path, prefixLength )
        handler( resp, req, suffix )
    }
}


func pathStripPrefix( path string, num int ) (string) {
    pathparts := strings.Split( path, "/" )
    if num >= len( pathparts ) { return "" }

    return strings.Join( pathparts[ num: ], "/" )
}

/*
func pathSuffix( path string, prefix string ) {
    preparts := strings.Split( prefix, "/" )
    pathparts := strings.Split( path, "/" )

}
*/