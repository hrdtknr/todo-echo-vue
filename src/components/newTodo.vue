<template>
  <div>
    <form class="newTodoForm">
      <input v-model="newName" placeholder="new name" />
      <input v-model="newTodo" placeholder="new todo" />
      <button v-on:click="createTodo()">新規登録</button>
    </form>
  </div>
</template>
<script>
import axios from "axios";
export default {
  data() {
    return {
      newName: "",
      newTodo: "",
    };
  },
  methods: {
    createTodo: function() {
      if (!this.newName && !this.newTodo) {
        alert("NameかTodoのどちらかは入力してください");
      } else {
        axios.post("http://localhost:8000/todoList", {
          name: this.newName,
          todo: this.newTodo,
        });
      }
    },
  },
};
</script>
<style scoped>
.newTodoForm {
  background-color: blueviolet;
  padding: 10px;
  width: 100%;
}
</style>
