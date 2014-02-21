$(function() {

    var pwd = [ '' ];

    mainEvents();

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
            .addClass( 'thumbnail' )
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

    function SlideShow( imageUrls ) {
        var self = this;
        this.imageUrls = imageUrls;
        this.index = 0;

        this.overlay = $('<div class="slideshow-overlay" />').appendTo( 'body' );
        this.container = $('<div class="slideshow" />').appendTo( 'body' );
        this.images = $('<ul/>').appendTo( this.container );

        $('<a href="#prev" class="slide-prev">Prev</a>')
            .click( function( ev ) {
                ev.preventDefault();
                self.prev();
            })
            .appendTo( this.container );

        $('<a href="#next" class="slide-next">Next</a>')
            .click( function( ev ) {
                ev.preventDefault();
                self.next();
            })
            .appendTo( this.container );

        $('<a href="#close" class="slide-close">Close</a>')
            .click( function( ev ) {
                ev.preventDefault();
                self.destroy();
            })
            .appendTo( this.container );


        for ( var i = 0; i < imageUrls.length; i++ ) {
            $('<li/>')
                .appendTo( this.images );
        }

        this.container.on( 'click', function( ev ) {
            ev.stopPropagation();
        });

        /*
        this.overlay.on( 'click', function( ev ) {
            self.destroy();
        });
        */
    }

    SlideShow.prototype = {

        destroy: function() {
            this.container.remove();
            this.overlay.remove();
        },

        showIndex: function( i ) {
            var li = this.images
                .children()
                .hide()
                .eq( i );

            var self = this;

            if ( !li.length )
                return;

            this.index = i;

            if ( !li.children('img').length ) {
                this.imgTag( i ).appendTo( li );
            }

            li.show();
            /*
            if ( li.children('img').width() ) {
                this.fit( li );
            } else {
                li.bind( 'load', function() {
                    self.fit( li );
                });
            }
            */
        },

        fit: function( li ) {
            var img = li.children('img');
            return img;
            //TODO: maybe fancy fitting?  seems css was able to handle it
        },

        next: function() {
            if ( (this.index + 1) < this.imageUrls.length ) {
                this.showIndex( this.index + 1 );
            }
        },

        prev: function() {
            if ( this.index > 0 ) {
                this.showIndex( this.index - 1 );
            }
        },

        imgTag: function( i ) {
            return $('<img />')
                .attr( 'src', this.imageUrls[i] );
        }
    };

    function mainEvents() {
        $(window).bind( 'hashchange', listFromHash );

        $('html').on( 'click', 'img.thumbnail', function( ev ) {
            ev.preventDefault();
            var images = $(this).closest( 'ul' ).find( 'a' ).map( function() {
                return $(this).attr( 'href' );
            });
            var slideshow = new SlideShow( images );
            slideshow.showIndex( $(this).closest( 'li' ).prevAll().length );
        });
    }
});
