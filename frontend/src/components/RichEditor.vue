<template>
  <div>
    <vue-editor v-model="content" :id="id" @input="textChange" @image-added="handleImageAdded" useCustomImageHandler></vue-editor>
  </div>
</template>

<script>
  import { VueEditor } from 'vue2-editor'
  
  export default {
    components: {
      VueEditor,
    },
    props: {
      id: {
        type: String,
        required: true
      },
      value: {
        type: String,
        required: false
      }
    },
    data() {
      return {
        content: ''
      }
    },
    watch: {
      value() {
        this.content = this.value
      }
    },
    methods: {
      textChange() {
        this.$emit('input', this.content)
      },
      handleImageAdded: function(file, Editor, cursorLocation, resetUploader) {
        //Editor.insertEmbed(cursorLocation, 'image', 'http://localhost:8080/images/45M_AN_O.jpg');
        var formData = new FormData();
        formData.append("file", file);
        
        var fileName = file.name.replace(/\.[^/.]+$/, "")

        this.axios.post('/images/' + fileName, formData, {headers: {'Content-Type': 'multipart/form-data'}}).then(function(result){
          let url = result.data.url; // Get url from response
          Editor.insertEmbed(cursorLocation, "image", url);
          resetUploader();
        }).catch(function(error){
          // eslint-disable-next-line
          console.log('upload error ' + error);
        });
      }
    },
    mounted () {
      // eslint-disable-next-line
      console.log('initialized RichEditor with id: ' + this.id)
      this.content = this.value
    }
  }
</script>

<style>

</style>