$(function() {

    var pwd = [ '' ];

    $(window).bind( 'hashchange', listFromHash );

    listFromHash();

    function listFromHash() {
        if ( window.location.hash.length ) {
            pwd = window.location.hash.split( '/' );
            pwd[0] = '';
            if ( !pwd[ pwd.length - 1 ].length )
                pwd.pop();  //remove element for trailing slash
        } else {
            pwd = [ '' ];
        }

        listDir( pwd );
    }

    function listDir( dir ) {
        dir = dir.join( '/' );
        $.getJSON( 'index.json', {root: dir}, function( data ) {
            if ( data.files ) {
                showIndex( data.files );
            }
        });
    }

    function showIndex( files ) {
        var thumbs = $('#index').empty(),
            subs = $('#folderlist').empty(),
            bread = $('#breadcrumb').empty(),
            i;

        for ( i = 0; i < files.length; i++ ) {
            if ( files[i].match( /\.(png|jpg|gif)/i ) ) {
                $('<li/>')
                    .append( thumbnail( files[i] ) )
                    .appendTo( thumbs );
            } else if ( files[i].match( /\/$/ ) ) {
                $('<li/>')
                    .append( subfolder( files[i] ) )
                    .appendTo( subs );
            }
        }

        for ( i = 0; i < pwd.length; i++ ) {
            $('<li/>')
                .append( subfolder( pwd.slice(0, i + 1).join( '/' ) ) )
                .appendTo( bread );
        }
    }

    function thumbnail( path ) {
        var link = $('<a/>')
            .attr( 'href', 'view/' + path );

        $('<img/>')
            .attr( 'src', 'thumb/' + path )
            .appendTo( link );

        return link;
    }

    function subfolder( path ) {
        if ( path.substr( -1 ) === '/' )
            path = path.substr( 0, path.length - 1 );

        return $('<a />')
            .attr( 'href', '#' + path )
            .text( path ? path.split( '/' ).pop() : 'Root' )
            .click( function( ev ) {
                pwd = path.split( '/' )
                /*
                if ( !pwd[ pwd.length - 1 ].length )
                    pwd.pop();  //remove element for trailing slash
                */
                //pwd.push( path );
                listDir( pwd );
            });
    }
});
