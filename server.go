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
    "container/heap"
    "encoding/json"
    "regexp"
    "path"
    "runtime"
    "strconv"
)

var config *ServerConfig
//var icache map[string]*bytes.Reader
var icache *ImageCache = NewImageCache()


func main() {
    if runtime.NumCPU() > 2 {
        runtime.GOMAXPROCS( runtime.NumCPU() -1 )
    }

    var err error
    config, err = loadImageCacheConfig()

    if err != nil {
        fmt.Println("Error parsing config: " + err.Error())
        return
    }

    fmt.Println(config)

	port := strconv.Itoa(config.Port)
	if len( os.Args ) > 1 {
		port = os.Args[1]
	}

	http.HandleFunc( "/", func( resp http.ResponseWriter, req *http.Request ) {
        //vals := req.URL.Path
        io.WriteString( resp, "Hello world!\n" + req.URL.Path )
    })

    http.HandleFunc( "/index.json", func(r http.ResponseWriter, req *http.Request) {
        indexJson(config, r, req)
    })

    relHandler( "/thumb/", serveThumb )

	listenOn := "localhost:" + port
	fmt.Println( "listening on " + listenOn )

	http.ListenAndServe(listenOn, nil )
}

func indexJson( config *ServerConfig, resp http.ResponseWriter, req *http.Request ) {
    var outerr error
    fmt.Println(config)
    defer func() {
        if recover() != nil {
            if outerr == nil {
                http.Error( resp, "Internal Server Error", 500 )
            } else {
                http.Error( resp, outerr.Error(), 500 )
            }
        }
    }()

    root := req.URL.Query().Get( "root" )
    if root != "" {
        root = path.Clean( root )
    }

    //TODO: make dir scanning possible.
    dir, err := os.Open( config.Root + "/" + root )
    defer dir.Close()

    if err != nil {
        log.Panic( "ERROR: Could not read directory" )
    }

    //TODO: Filter out non-images
    filestats, err := dir.Readdir( -1 )
    if err != nil {
        log.Panic( "ERROR: Could not read directory" )
    }


    files := make([]string, 0, len(filestats))
    for _, file := range filestats {
        if file.IsDir() && file.Name()[0] != '.' {
            files = append( files, root + "/" + file.Name() + "/" )
        } else if isImage( file ) {
            files = append( files, root + "/" + file.Name() )
        }
    }

    out := map[string]interface{}{
        "status": "OK",
        "files": files,
    }

    if bytes, err := json.Marshal( out ); err == nil {
        resp.Write( bytes )
    } else {
        io.WriteString( resp, "[]" )
    }
}

var imgRE *regexp.Regexp = nil
func isImage( f os.FileInfo ) (bool) {
    if imgRE == nil {
        imgRE = regexp.MustCompile( `(?i)\.(png|jpg|gif)$` )
    }

    return imgRE.MatchString( f.Name() )
}

func serveThumb( resp http.ResponseWriter, req *http.Request, path string ) {

    var outerr error

    log.Println( "Requested image", path )
    //log.Println( "Oldest: ", icache.Top() )
    //log.Println( "heap: ", icache.GetPaths() )

    defer func() {
        if recover() != nil {
            if outerr == nil {
                http.Error( resp, "Internal Server Error", 500 )
            } else {
                http.Error( resp, outerr.Error(), 500 )
            }
        }
    }()

    /*
    if icache == nil {
        //icache = make( map[string]*bytes.Reader, 100 )
        icache = &ImageCache{}
    }
    */

    cached, ok := icache.Find( path )

    if ok {
        icache.Update( path )
    } else {
        log.Println( "Generating thumbnail" )
        fd, err := os.Open( config.Root + "/" + path )
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

        heap.Push( icache, NewCacheItem( path, cached ) )

        if icache.Len() > config.MaxImages {
            log.Println( "Dropping oldest cache: ", heap.Pop( icache ).(*CacheItem).path )
        }

        //icache[ path ] = cached
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
