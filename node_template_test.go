package node_template;

import (
        "testing"
        "strings"
	"bytes"
	"container/list"
);


func TestParseRender(t *testing.T ){
	temp, err := Parse( strings.NewReader("<html><head></head><body><div>foo</div></body></html>") );
	if( err != nil ){
		t.Error(err);
		return;
	}
	var out bytes.Buffer 
	temp.Render( &out );
	if( out.String() != "<html><head></head><body><div>foo</div></body></html>"){
		t.Error("Cant get identical output from Parse -> Render" + out.String());
	}
}

func TestNodeTemplateFromFile( t * testing.T ){
	temp, err := NodeTemplateFromFile( "tdata/t.html");
	if( err != nil ){
		t.Error(err);
		return;
	}
	var out bytes.Buffer 
	temp.Render( &out );
	s := out.String();
	s = strings.Replace( s, "\n", "", -1);
	s = strings.Replace( s, " ", "", -1);
	if( s != "<html><head></head><body><div>foo</div></body></html>"){
		t.Error("Cant get identical output from NodeTemplateFromFile -> Render:`" + s + "`");
	}

}
func TestReplaceText( t * testing.T ){	
	temp, _ := Parse( strings.NewReader("<html><head></head><body><div id=foo>foop</div></body></html>") );
	el, _ := temp.FindFirst("#foo");
	el.ReplaceContentText("bar");
	var out bytes.Buffer 
	temp.Render( &out );
	s := out.String();
	if( strings.Index(s, "bar") == -1 ){
		t.Error("Cant find replace text in string :`" + s + "`");
	}
	if( strings.Index(s, "foop") != -1 ){
		t.Error("Replaced text still in string :`" + s + "`");
	}

}
func TestFindFirst( t * testing.T ){
	temp, _ := Parse( strings.NewReader("<html><head></head><body><div id=foo>foop</div></body></html>") );
	el, _ := temp.FindFirst("#foo");
	if( el == nil ){
		t.Error("Can't find id in template");
	}
	el2, _ := temp.FindFirst("#bar");
	if( el2 != nil ){
		t.Error("find ID not in template ");
	}
}

func TestFind( t * testing.T ){
	temp, _ := Parse( strings.NewReader("<html><head></head><body><div id=foo>foop</div><span class='chi'>one</span><span class='chi'>two</span></body></html>") );
	nodes, _ := temp.Find("#foo");
	if( nodes.Len() != 1 ){
		t.Error("Can't find id in template");
	}
	nodes2, _ := temp.Find("#bar");
	if( nodes2.Len() != 0 ){
		t.Error("find ID not in template ");
	}
	nodes3, _ := temp.Find(".chi");
	if( nodes3.Len() != 2 ){
		t.Error("find multiple classes doesn't return expected amount");
	}
	if( nodes3.Get(0).FirstChild.Data != "one" ){
		t.Error("First class node data not as expected: " + nodes3.Get(0).Data);
	}
	if( nodes3.Get(1).FirstChild.Data != "two" ){
		t.Error("First class node data not as expected: " + nodes3.Get(0).Data);
	}
}
func TestRepeatNode( t * testing.T ){
	template, _ := Parse( strings.NewReader("<html><head></head><body><div class='person'>foop</div></body></html>") );
	
	var people list.List
	people.PushBack("Tom");
	people.PushBack("Richard");
	people.PushBack("Harry");
        person, _ := template.FindFirst(".person");
        if( person != nil ){
		person.RepeatNode( &people, func( node * NodeTemplate, e * list.Element  ){
			node.ReplaceContentText( e.Value.(string) );
		});
        }
	results, _ := template.Find(".person");
	
	if( results.Len() != 3  ){
		t.Error("Dont get expected amount of elements with .person class");
	}
	if( results.Get(0).FirstChild.Data != "Tom" ){
		t.Error("Dont get Tom in 1");
	}
	if( results.Get(1).FirstChild.Data != "Richard" ){
		t.Error("Dont get Richard in 2");
	}
	if( results.Get(2).FirstChild.Data != "Harry" ){
		t.Error("Dont get Harry in 3");
	}

}
