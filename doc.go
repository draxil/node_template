/*
node_template extends go.net/html & cascadia to make it as easy as possible to
process pure HTML templates. Using jquery like search syntax you can find the
elements you wish to populate by looking for the attributes you require.
This approach allows HTML templates which are pure HTML.


Example:


 func main_page( w http.ResponseWriter, r * http.Request ){
	template, err := node_template.NodeTemplateFromFile( "temp/main.html");

	if( err != nil ){
		log.Println( err );
		return
	}

       	title_el, err := template.FindFirst(`#title`);
        if( title_el != nil){
           title_el.ReplaceContentText("Billy & Jane");
        }

        names, err := template.Find(`.name`);
        if( names != nil ){
            names.ReplaceContentText("tom");
        }

	var people list.List
	people.PushBack("Tom");
	people.PushBack("Richard");
	people.PushBack("Harry");
        person, _ := template.FindFirst(".person");
        if( person != nil ){
		person.RepeatNode( &people, func( node * node_template.NodeTemplate, e * list.Element  ){
			node.ReplaceContentText( e.Value.(string) );
		});
        }

	template.Render( w );
 }


*/