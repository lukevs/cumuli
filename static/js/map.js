// JS for creating the soundcloud map with D3

var width = window.innerWidth/2, 
    height = window.innerHeight;

var linkDistance = 100;
var spinLength = 15;

// Handle mobile
if (width < (768 / 1.9)) {
    height = height/2;
    linkDistance = 50;
    spinLength = 5
}

var force = d3.layout.force()
    .charge(-130)
    .linkDistance(linkDistance)
    .size([width, height])
    .friction(0.8);

var svg = d3.select("#graph").append("svg")
    .attr("width", width)
    .attr("height", height);

// Get the JSON path from an attribute to the source tag
var scriptParam = document.getElementById('soundcloud-map')
var jsonPath = scriptParam.getAttribute('jsonPath')

// Loading gif
var spinner = new Spinner({
  lines: 9,
  length: spinLength,
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

            // Handle mobile
            if (width < (768 / 1.9)) {
                if (d.group != 1) {
                    return d.weight * 2; 
                }
                return 7;
            }

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

    // Hover
    $('svg circle').tipsy({ 
        trigger: 'click',
        gravity: 'e',
        html: true, 
        title: function() {
          var d = this.__data__;
          return d.name; 
        }
    });

    var r = 10;

    force.on("tick", function() {
        node.attr("cx", function(d) { return d.x = Math.max(r, Math.min(width - r, d.x)); })
            .attr("cy", function(d) { return d.y = Math.max(r, Math.min(height - r, d.y)); });

        link.attr("x1", function(d) { return d.source.x; })
            .attr("y1", function(d) { return d.source.y; })
            .attr("x2", function(d) { return d.target.x; })
            .attr("y2", function(d) { return d.target.y; });

    });
});