// JS for creating the soundcloud map with D3

var width = window.innerWidth/1.7, 
    height = window.innerHeight;

var color = d3.scale.category20()

var force = d3.layout.force()
    .charge(-120)
    .linkDistance(100)
    .size([width, height]);

var svg = d3.select("#graph").append("svg")
    .attr("width", width)
    .attr("height", height);

// Get the JSON path from an attribute to the source tag
var scriptParam = document.getElementById('soundcloud-map')
var jsonPath = scriptParam.getAttribute('jsonPath')

// Loading gif
var spinner = new Spinner({
  lines: 9,
  length: 15,
  width: 5,
  radius: 20,
  color: '#4A93A2'
}).spin(document.getElementById("spinner-box"));

// D3
d3.json(jsonPath, function(error, graph) {
    force
        .nodes(graph.nodes)
        .links(graph.links)
        .start()

    var link = svg.selectAll(".link")
        .data(graph.links)
        .enter().append("line")
        .attr("class", "link");

    var node = svg.selectAll(".node")
        .data(graph.nodes)
        .enter().append("circle")
        .attr("class", "node")
        .attr("r", function(d) { 
            if (d.group != 1) {
                return d.weight * 3; 
            }
            return 10;
        })
        .style("fill", function(d) {
            if (d.group == 1) {
                return "#FA6900"
            }
            return "#4A93A2"
        })
        .call(force.drag);

    node.append("title")
        .text(function(d) { return d.name; });

    // Remove the loading gif
    spinner.stop()

    force.on("tick", function() {
        link.attr("x1", function(d) { return d.source.x; })
            .attr("y1", function(d) { return d.source.y; })
            .attr("x2", function(d) { return d.target.x; })
            .attr("y2", function(d) { return d.target.y; });

        node.attr("cx", function(d) { return d.x; })
            .attr("cy", function(d) { return d.y; });
    });
});