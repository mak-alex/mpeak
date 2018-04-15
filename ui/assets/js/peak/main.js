var menu_data = [
  { id:"mountains", icon: "", value: "Mountains", hotkey: "enter+shift"},
  { id:"contact", icon:"", value:"Contact Us" }
]; 

var about = [
  { template:"About author and soft", type:"header" },
  { view:"text", value:'Alex M.A.K.', label:"Author", readonly:true},
  { view:"text", value:'alex-m.a.k@yandex.kz', label:"E-Mail", readonly:true},
  { view:"text", value:"mPeak", label:"Program", readonly:true },
  { view:"text", value:"0.1.0", label:"Version", readonly:true },
  { view:"text", value:"https://bitbucket.org/enlab/mpeak", label:"Link", readonly:true },
  { view:"text", value:"15.04.18", label:"Created", readonly:true },
];

var add_mountain = {
  view: "form",
  id: "addMountain",
  width: 350,
  elements: [{
    view: "text",
    label: "Title",
    name: "title"
  }, {
    view: "text",
    label: "Latitude",
    name: "latitude",
  }, {
    view: "text",
    label: "Longtitude",
    name: "longtitude",
  }, {
    view: "text",
    label: "Height",
    name: "height",
  }, {
    view: "text",
    label: "WebID",
    name: "web_id",
  }, {
    view: "text",
    label: "TypeLink",
    name: "type_link",
  }, {
    view: "button",
    value: "Add",
    width: 150,
    align: "center",
    click:function(){
      if (this.getParentView().validate()){ //validate form
        var t = $$("addMountain").getValues();
        webix.ajax().headers({ "Content-type":"application/json" }).put("/api/v1/mountains", t, 
        function(text, xml, xhr){
          console.log(text);
        });
        this.getTopParentView().hide(); //hide window
      }
      else
        webix.message({ type:"error", text:"Form data is invalid" });
    }
  }]
};

webix.ui({
  view:"window",
  id:"add_mountain_window",
  width:300,
  position:"center",
  modal:true,
  move: true,
  head: {
     view: "toolbar", margin: -4, cols: [
         { view: "label", label: "New Mountain" },
         {
             view: "icon", icon: "question-circle",
             click: "webix.message('About pressed')"
         },
         {
             view: "icon", icon: "times-circle",
             click: '$$("add_mountain_window").hide();'
         }
     ]
  },
  body:add_mountain
});

function showForm(winId, node){
  $$(winId).getBody().clear();
  $$(winId).show(node);
  $$(winId).getBody().focus();
};

function submit(){
  webix.message(JSON.stringify($$("addMountain").getValues(), null, 2));
};

var header = {
  view:"toolbar",
  elements:[
    { view:"label",  label: "Peak Simple Client"},
    { view:"text", id:"grouplist_input", placeholder:"Поиск ..."},
    { view:"button", width:50, value: 'Add', click:function(){ showForm("add_mountain_window") }},
  ]
};

var mountains= { id:"mountains", css:"preview-box", rows:[
  {
    id:"mountains_list",
    view:"datatable",
    columns:[
      { id:"id", header:"", css:{"text-align":"center"}, width:50, hidden: true},
      { id:"title",	header:"Name", 	width:200, fillspace: 6 },
      { id:"latitude",	header:"Latitude", 	width:200, fillspace: 4 },
      { id:"longtitude",	header:"Longtitude", 	width:200, fillspace: 4 },
      { id:"height",	header:"Height", 	width:200, fillspace: 4 },
      { id:"web_id",	header:"WebID", 	width:200, fillspace: 4 },
      { id:"type_link",	header:"TypeLink", 	width:200, fillspace: 4 },
      { id:"trash", header:"", template:"{common.trashIcon()}", fillspace: 1}
    ],
    datafetch: 20,
    loadahead: 20,
    url:"/api/v1/mountains",
    onClick:{
      "fa-trash":function(event, id, node){
        webix.message("Delete row: "+id);
        this.remove(id);
        webix.ajax().headers({ "content-type":"application/json" }).del("/api/v1/mountains/"+id, {}, function(text, xml, xhr){
          console.log(text)
        });
        return false; // here it blocks the default behavior
      }
    },
    on: {
      onDataRequest: function (start, count) {
        onBeforeFilterCount++;
        webix.ajax().get("/api/v1/mountains?page="+onBeforeFilterCount+"&per_page="+count).then(function(data){
          console.log(data.json());
          this.clearAll();
          this.parse(data.json());
        })
        return false;
      },
    }
  }
]}

var contacts = { id:"contact", cols: [
  { view:"form", elements:about, scroll:true },
]}

var onBeforeFilterCount = 0;
webix.ready(function(){
    if (webix.CustomScroll && !webix.env.touch)
        webix.CustomScroll.init();

    webix.ui({
      rows:[
        header,
        { animate:{type:"slide",subtype:"out", direction:"bottom"}, cells:[
          mountains,
          contacts
        ]},
        { view:"tabbar", type:"bottom", id:"tab", height:40, options:menu_data, multiview:true }
      ]
    });

    webix.ui.fullScreen();
});
