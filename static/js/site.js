$(document).ready(function() {

    // Variables
    var $nav = $('.navbar');
    if ($nav !== undefined && $nav !== null && $nav.top !== undefined && $nav.top !== null) {
      var $navOffsetTop = $nav.offset().top;
    }
    var $body = $('body'),
        $window = $(window),
        $popoverLink = $('[data-popover]'),
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
      if ($nav !== undefined) {
        navOffsetTop = $nav.offset().top
        onScroll()
      }
    }
  
    function onScroll() {
      if(typeof navOffsetTop !== "undefined" && navOffsetTop < $window.scrollTop() && !$body.hasClass('has-docked-nav')) {
        $body.addClass('has-docked-nav')
      }
      if(typeof navOffsetTop !== "undefined" && navOffsetTop > $window.scrollTop() && $body.hasClass('has-docked-nav')) {
        $body.removeClass('has-docked-nav')
      }
    }
  
    function escapeHtml(string) {
      return String(string).replace(/[&<>"'\/]/g, function (s) {
        return entityMap[s];
      });
    }
  
    $("#pagesdelete").click(function() {

      var pagesToDeleteUUIDs = [];

      $("#page-list tr").each(function(){
        collectAllCheckedBoxIDs(this, pagesToDeleteUUIDs);
      })

      if (pagesToDeleteUUIDs.length > 0) {

        if (confirm("Delete " + String(pagesToDeleteUUIDs.length) + " page" + ((pagesToDeleteUUIDs.length > 1) ? "s?" : "?"))) {
          console.log()
          var form = document.createElement("form");
          form.setAttribute("id", "deleteform");
          form.setAttribute("method", "POST");
          form.setAttribute("action", window.location.pathname + "/delete");
  
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

    $("#usersdelete").click(function() {

      var usersToDeleteUUIDs = [];

      $("#user-list tr").each(function(){
        collectAllCheckedBoxIDs(this, usersToDeleteUUIDs);
      })

      if (usersToDeleteUUIDs.length > 0) {
        if (confirm("Delete " + String(usersToDeleteUUIDs.length) + " user" + ((usersToDeleteUUIDs.length > 1) ? "s?" : "?"))) {
          var form = document.createElement("form");
          form.setAttribute("id", "deleteform");
          form.setAttribute("method", "POST");
          form.setAttribute("action", window.location.pathname + "/delete");

          form._submit_function_ = form.submit;

          for (var i = 0; i < usersToDeleteUUIDs.length; i++) {
            var hiddenField = document.createElement("input");
            hiddenField.setAttribute("type", "hidden");
            hiddenField.setAttribute("name", String(i));
            hiddenField.setAttribute("value", usersToDeleteUUIDs[i]);
            form.appendChild(hiddenField);
          }
          document.body.appendChild(form);
          form._submit_function_();
        }
      }
    })

    $("#adduserstogroup").click(function() {

      var usesrToAddUUIDs = [];

      $("#users-not-in-group-list tr").each(function(){
        collectAllCheckedBoxIDs(this, usesrToAddUUIDs);
      })

      if (usesrToAddUUIDs.length > 0) {
        if (confirm("Add " + String(usesrToAddUUIDs.length) + " user" + ((usesrToAddUUIDs.length > 1) ? "s" : "") + " to group?")) {
          var form = document.createElement("form");
          form.setAttribute("id", "addusertogroupform");
          form.setAttribute("method", "POST");
          form.setAttribute("action", window.location.pathname + "/add");

          form._submit_function_ = form.submit;

          for (var i = 0; i < usesrToAddUUIDs.length; i++) {
            var hiddenField = document.createElement("input");
            hiddenField.setAttribute("type", "hidden")
            hiddenField.setAttribute("name", String(i));
            hiddenField.setAttribute("value", usesrToAddUUIDs[i]);
            form.appendChild(hiddenField);
          }
          document.body.appendChild(form);
          form._submit_function_();
        }
      }
    })

    $("#selectallpages").change(function() {
      var selectAll = this.checked;
      $("#page-list tr").each(function(){
        selectAllCheckboxes(this, selectAll);
      })
    });

    $("#selectallusers").change(function() {
      var selectAll = this.checked;
      $("#user-list tr").each(function(){
        selectAllCheckboxes(this, selectAll);
      })
    });

    $("#selectallusersingroup").change(function() {
      var selectAll = this.checked;
      $("#users-in-group-list tr").each(function(){
        selectAllCheckboxes(this, selectAll)
      })
    });

    $("#selectallusersnotingroup").change(function() {
      var selectAll = this.checked;
      $("#users-not-in-group-list tr").each(function(){
        selectAllCheckboxes(this, selectAll)
      })
    });

    function selectAllCheckboxes(row, selectAll) {
      $(row).find("td").each(function(){
        $(this).find("input").each(function(){
          if ($(this).attr("type") == "checkbox") {
            $(this).prop("checked", selectAll);
          }
        })
      })
    }

    function collectAllCheckedBoxIDs(row, collection) {
      $(row).find("td").each(function(){
        var tableCellId = $(this).attr("id");
        if (tableCellId) {
          $(this).find("input").each(function(){
            if ($(this).attr("type") == "checkbox") {
              if ($(this).is(":checked")) {
                collection.push(tableCellId);
              }
            }
          })
        }
      })
    }

    $('#newrootform').submit(function(){
      return validatePasswords();
    });

    function validatePasswords() {
      if ($("#newpass").val().length > 0 && $("#repnewpass").val().length > 0) {
        if ($("#newpass").val() !== $("#repnewpass").val()) {
          $('#repnewpass').get(0).setCustomValidity("Passwords don't match");
          return false;
        }
        $('#repnewpass').get(0).setCustomValidity("");
        return true;
      }
      return false;
    }

    $('#repnewpass').keyup(validatePasswords);
  
    init();
  
  });