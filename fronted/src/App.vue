<template>
  <div class="app">
    <newTodo></newTodo>
    <todoTable></todoTable>
  </div>
</template>

<script>
import todoTable from "./components/todoTable";
import newTodo from "./components/newTodo";
export default {
  components: {
    todoTable,
    newTodo,
  },
};
</script>
<style scoped>
.app {
  background-color: greenyellow;
  padding: 10px;
  width: 100%;
}
</style>
