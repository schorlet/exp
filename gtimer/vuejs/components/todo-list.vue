<template>
	<section
			class="todo-list"
			v-show="todos.length"
	>
		<article
			is="todo-item"
			v-for="(todo, index) in todos"
				:todo="todo"
				:key="todo.id"
				:index="index"
			@toggle="onToggle"
			@update="onUpdate"
			@remove="onRemove"
		></todo-item>
	</section>
</template>

<script>
module.exports = {
	name: 'TodoList',
	components: {
		TodoItem: httpVueLoader('./todo-item.vue')
	},
	props: {
		todos: {
			type: Array,
			required: true
		}
	},
	methods: {
		onToggle: function(id) {
			this.$emit('toggle', id);
		},
		onUpdate: function(todo) {
			this.$emit('update', todo);
		},
		onRemove: function(id) {
			this.$emit('remove', id);
		}
	}
}
</script>

<style scoped>
	.todo-list article {
		margin: 6px;
		border: 1px solid orange;
	}
</style>
