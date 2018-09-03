$(document).ready(function() {

    // Variables
    var $nav = $('.navbar'),
        $body = $('body'),
        $window = $(window),
        $popoverLink = $('[data-popover]'),
        navOffsetTop = $nav.offset().top,
        $document = $(document),
        entityMap = {
          "&": "&amp;",
          "<": "&lt;",
          ">": "&gt;",
          '"': '&quot;',
          "'": '&#39;',
          "/": '&#x2F;'
        }
  
    function init() {
      $window.on('scroll', onScroll)
      $window.on('resize', resize)
      $popoverLink.on('click', openPopover)
      $document.on('click', closePopover)
      $('a[href^="#"]').on('click', smoothScroll)
    }
  
    function smoothScroll(e) {
      e.preventDefault();
      $(document).off("scroll");
      var target = this.hash,
          menu = target;
      $target = $(target);
      $('html, body').stop().animate({
          'scrollTop': $target.offset().top-40
      }, 0, 'swing', function () {
          window.location.hash = target;
          $(document).on("scroll", onScroll);
      });
    }
  
    function openPopover(e) {
      e.preventDefault()
      closePopover();
      var popover = $($(this).data('popover'));
      popover.toggleClass('open')
      e.stopImmediatePropagation();
    }
  
    function closePopover(e) {
      if($('.popover.open').length > 0) {
        $('.popover').removeClass('open')
      }
    }
  
    $("#button").click(function() {
      $('html, body').animate({
          scrollTop: $("#elementtoScrollToID").offset().top
      }, 2000);
  });
  
    function resize() {
      $body.removeClass('has-docked-nav')
      navOffsetTop = $nav.offset().top
      onScroll()
    }
  
    function onScroll() {
      if(navOffsetTop < $window.scrollTop() && !$body.hasClass('has-docked-nav')) {
        $body.addClass('has-docked-nav')
      }
      if(navOffsetTop > $window.scrollTop() && $body.hasClass('has-docked-nav')) {
        $body.removeClass('has-docked-nav')
      }
    }
  
    function escapeHtml(string) {
      return String(string).replace(/[&<>"'\/]/g, function (s) {
        return entityMap[s];
      });
    }
  
    $("#delete").click(function() {

      var pagesToDeleteUUIDs = [];

      $("#page-list tr").each(function(){
        $(this).find("td").each(function(){
          var tableCellId = $(this).attr("id");
          if (tableCellId) {
            $(this).find("input").each(function(){
              if ($(this).attr("type") == "checkbox") {
                if ($(this).is(":checked")) {
                  pagesToDeleteUUIDs.push(tableCellId);
                }
              }
            })
          }
        })
      })

      if (pagesToDeleteUUIDs.length > 0) {

        if (confirm("Delete " + String(pagesToDeleteUUIDs.length) + " page" + ((pagesToDeleteUUIDs.length > 1) ? "s?" : "?") )) {
          var form = document.createElement("form");
          form.setAttribute("id", "deleteform");
          form.setAttribute("method", "POST");
          form.setAttribute("action", "/admin/pages/delete");
  
          form._submit_function_ = form.submit;
  
          for (var i = 0; i < pagesToDeleteUUIDs.length; i++) {
            var hiddenField = document.createElement("input");
            hiddenField.setAttribute("type", "hidden");
            hiddenField.setAttribute("name", String(i));
            hiddenField.setAttribute("value", pagesToDeleteUUIDs[i]);
            form.appendChild(hiddenField);
          }
          document.body.appendChild(form);
          form._submit_function_();
        }

      }
    });
  
    init();
  
  });