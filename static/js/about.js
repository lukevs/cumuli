$('.clickable').tipsy({ 
    trigger: 'click',
    gravity: 'e',
    html: true, 
    title: function() {
        return this.getAttribute('original-title')
    }
});