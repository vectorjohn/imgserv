$(function() {
    $.getJSON( 'index.json', function( data ) {
        if ( data.files ) {
            showIndex( data.files );
        }
    });


    function showIndex( files ) {
        var ul = $('#index').empty();

        for ( var i = 0; i < files.length; i++ ) {
            if ( files[i].match( /\.(png|jpg|gif)/i ) ) {
                $('<li/>')
                    .append( thumbnail( files[i] ) )
                    .appendTo( ul );
            }
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
});
