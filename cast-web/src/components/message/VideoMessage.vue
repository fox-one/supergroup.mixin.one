<template>
  <vue-plyr
    ref="plyr"
    :options="options"
    class="video"
  >
    <video>
      <source
        :src="message.attachment.view_url"
        :type="message.attachment.mime_type"
      >
    </video>
  </vue-plyr>
</template>
<script>
import plyrMixin from '@/mixins/plyr'

export default {
  name: "VideoMessage",
  mixins: [plyrMixin],
  props: {
    message: {
      type: Object,
      default: () => {}
    }
  },
  data () {
    return {
      options: {
        controls: ['play-large', 'play', 'progress', 'current-time', 'mute', 'fullscreen']
      }
    };
  }
}
</script>
<style lang="scss" scoped>
.video {
  max-width: 220px;
}
</style>