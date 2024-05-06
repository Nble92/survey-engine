<template>
  <v-col cols="12">
    <v-row no-gutters>
      <div class="pb-4" v-if="error" style="width: 100%">
        <v-alert
          v-model="error"
          type="error"
          transition="scale-transition"
        >{{ errorText }}</v-alert>
      </div>
      
      <v-dialog v-model="processingDialog" hide-overlay persistent width="300">
        <v-card color="primary" dark class="pt-3">
          <v-card-text>
            Processing...
            <v-progress-linear
              indeterminate
              color="white"
              class="mb-0"
            ></v-progress-linear>
          </v-card-text>
        </v-card>
      </v-dialog>
      
      <v-dialog v-model="imageDeleteDialog" persistent max-width="600px">
        <v-card>
          <v-card-title>
            <span class="headline">Delete Image</span>
          </v-card-title>
          <v-card-text class="text-xs-center">
            Are you sure you want to delete the image '<strong>{{ deleteImageName }}</strong>'?
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="imageDeleteDialog = false">No</v-btn>
            <v-btn color="blue darken-1" text @click="deleteImageConfirm()">Yes, Delete it</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    
      <v-dialog v-model="imageUploadDialog" persistent max-width="600px">
        <v-card>
          <v-card-title>
            <span class="headline">Upload Image</span>
          </v-card-title>
          <v-card-text class="pa-4">
            <v-form>
              <!-- TODO fix this uglyness -->
              <div class="v-text-field__slot mb-3">
                <input type="file" accept="image/*" class="" aria-label="Image" required="required" @change="fileChange($event.target.files)" ref="fileInput">
              </div>
              <v-text-field
                v-model="imageName"
                label="Image Name"
                required
              ></v-text-field> 
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="imageUploadDialog = false">Cancel</v-btn>
            <v-btn color="blue darken-1" text @click="uploadImage()">Upload</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      
      <v-card style="width: 100%">
        <v-toolbar class="blue-grey white--text">
          <v-toolbar-title class="headline font-weight-regular"><v-icon dark class="mr-2">collections</v-icon> Images</v-toolbar-title>
          <v-spacer></v-spacer>
          <v-btn @click="imageUploadDialog = true" outlined color="white">Upload Image</v-btn>
        </v-toolbar>
        <v-card-title>
          <v-text-field
            v-model="search"
            append-icon="search"
            label="Search"
            single-line
            hide-details
          ></v-text-field>
        </v-card-title>
        <v-data-table
          :headers="imageListHeaders"
          :items="images"
          class="elevation-1"
          :search="search"
          :footer-props="{
              itemsPerPageOptions: [20, 50, 100],
            }"
        >
        
          <template v-slot:item.imageName="{ item }">
            <span class="font-weight-medium">{{ item.imageName }}</span>
          </template>
          
          <template v-slot:item.url="{ item }">
            <a :href="item.url" target="_blank">{{ baseURL }}{{ item.url }}</a>
          </template>
          
          <template v-slot:item.actions="{ item }">
            <v-icon
              small
              @click="deleteImage(item.imageName)"
            >
              delete
            </v-icon>
          </template>
        </v-data-table>
      </v-card>
    </v-row>
  </v-col>
</template>


<script>
  import { mapState } from 'vuex'

  export default {
    data () {
      return {
        valid: false,
        error: false,
        errorText: 'an unknown error occurred',
        
        processingDialog: false,
        search: '',
        
        imageListHeaders: [
          { text: 'Image Name', value: 'imageName' },
          { text: 'Image URL', value: 'url' },
          { text: 'Actions', value: 'actions', sortable: false, align: 'center' }
        ],
        
        imageUploadDialog: false,
        imageName: '',
        image: false,
        
        // deletion stuff
        imageDeleteDialog: false,
        deleteImageName: '',
      }
    },
    computed: {
      baseURL() {
        return location.protocol + '//' + location.hostname + ':' + location.port
      },
      ...mapState('images', ['images']),
    },
    mounted() {
      var that = this
      this.$store.dispatch('images/fetchImages').catch(function (error) {
        if (error.response && error.response.data) {
          that.errorText = error.response.data.message
        } else {
          that.errorText = error.message
        }
        that.error = true
      })
    },
    methods: {
      fileChange(fileList) {
        this.image = fileList[0]
        
        this.imageName = fileList[0].name
        
        this.imageName = this.imageName.replace(/\.[^/.]+$/, "")
      },
      uploadImage() {
        if (!this.imageName) {
          return
        }

        
        var that = this
        this.imageUploadDialog = false
        this.processingDialog = true
        
        setTimeout(function() {
          that.processingDialog = false
        }, 1000)
        
        let formData = new FormData();
        formData.append('file', this.image);
        
        this.axios.post('/images/' + this.imageName, formData, {headers: {'Content-Type': 'multipart/form-data'}}).then(function(){
          that.$store.dispatch('images/fetchImages')
          
          that.image = false
          that.imageName = ''
          that.$refs.fileInput.value = ''
        }).catch(function(error){
          if (error.response && error.response.data) {
            that.errorText = error.response.data.message
          } else if (error.message) {
            that.errorText = error.message
          } else {
            that.errorText = "Internal Error"
          }
          that.error = true
        });
      },
      deleteImage(name) {
        this.deleteImageName = name
        this.imageDeleteDialog = true
      },
      deleteImageConfirm() {
        var that = this
        this.processingDialog = true
        
        setTimeout(function() {
          that.processingDialog = false
        }, 1000)
        
        this.axios.delete('/images/' + this.deleteImageName)
        .then(() => {
          that.$store.dispatch('images/fetchImages')
        }).catch(function(error) {
          if (error.response && error.response.data) {
            that.errorText = error.response.data.message
          } else if (error.message) {
            that.errorText = error.message
          } else {
            that.errorText = "Internal Error"
          }
          that.error = true
        })
        
        this.deleteImageName = ''
        this.imageDeleteDialog = false
      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>