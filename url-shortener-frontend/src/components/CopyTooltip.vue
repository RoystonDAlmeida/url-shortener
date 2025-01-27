<template>
  <div class="tooltip" v-if="visible" :style="{ top: position.top, left: position.left }">
    {{ text }}
  </div>
</template>

<script>
export default {
  name: 'CopyTooltip',
  props: {
    text: {
      type: String,
      required: true
    },
    position: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      visible: false
    };
  },
  methods: {
    show() {
      this.visible = true;
      setTimeout(() => this.visible = false, 2000); // Hide after 2 seconds
    }
  }
};
</script>

<style scoped>
  .tooltip {
    position: absolute;
    background-color: black;
    color: white;
    padding: 5px;
    border-radius: 3px;
  }
</style>