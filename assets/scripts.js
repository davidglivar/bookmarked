var APP = (function () {
  var form = document.forms.bookmarked,
    submit = form.elements.commit,
    list = document.getElementsByTagName('ul')[0];
  
  var that = {

    destroyBookmark: function () {
      var destroyers = document.getElementsByClassName('destroyer');
      for (var i in destroyers) {
        destroyers[i].onclick = function (e) {
          e.preventDefault();
          var xhr = new XMLHttpRequest(),
            formdata = new FormData(),
            li = this.parentElement;

          xhr.open("DELETE", "/delete_bookmark", true);
          formdata.append("id", this.attributes['data-id'].value)
          xhr.send(formdata);

          xhr.onreadystatechange = function () {
            if (xhr.readyState === 4 && xhr.status === 200) {
              li.parentElement.removeChild(li);
            }
          }
        }
      }
    },

    saveBookmark: function() {
      submit.onclick = function (e) {
        e.preventDefault();
        var xhr = new XMLHttpRequest(),
          formdata = new FormData();

        xhr.open("POST", "/create_bookmark", true);
        formdata.append("name", form.name.value);
        formdata.append("url", form.url.value);
        xhr.send(formdata);

        xhr.onreadystatechange = function () {
          if (xhr.readyState === 4 && xhr.status === 200) {
            var li = document.createElement('li'),
              a = document.createElement('a'),
              txt = document.createTextNode(form.name.value);

            a.setAttribute('href', form.url.value);
            list.appendChild(li);
            li.appendChild(a);
            a.appendChild(txt);
            form.reset();
          }
        }
      }
    },

    init: function () {
      that.saveBookmark();
      that.destroyBookmark();
    }
  };
  return that;

}());

APP.init();
