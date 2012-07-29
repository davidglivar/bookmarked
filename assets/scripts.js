var APP = (function () {
  
  var that = {

    saveBookmark: function() {
      var form = document.forms.bookmarked,
        submit = form.elements.commit,
        list = document.getElementsByTagName('ul')[0];

      submit.onclick = function (e) {
        var xhr = new XMLHttpRequest(),
          formdata = new FormData();

        e.preventDefault();

        xhr.open("POST", "/create_bookmark", true);
        formdata.append("name", form.name.value);
        formdata.append("url", form.url.value);
        xhr.send(formdata);

        xhr.onreadystatechange = function () {
          if (xhr.readyState === 4 && xhr.status === 200) {
            console.log('success.');
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
    }
  };
  return that;

}());

APP.init();
