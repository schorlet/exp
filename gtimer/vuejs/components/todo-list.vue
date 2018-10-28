<template>
	<div
			class="todo-list"
			v-show="todos.length"
	>
		<todo-item
			v-for="(todo, index) in todosFiltered"
				:todo="todo"
				:key="todo.id"
				:index="index"
			@toggle="onToggle"
			@update="onUpdate"
			@remove="onRemove"
			:highlight="highlight"
			@drop="onDrop"
		></todo-item>

		<div class="filters">
			<a href="#/all"
				:class="{selected: filter == 'all'}"
			>All</a>
			<a href="#/active"
				:class="{selected: filter == 'active'}"
			>Active</a>
			<a href="#/completed"
				:class="{selected: filter == 'completed'}"
			>Completed</a>
		</div>
	</div>
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
		},
		highlight: {
			type: String,
			default: ''
		}
	},
	data: function() {
		return {
			filter: ''
		}
	},
	created: function() {
		window.addEventListener('hashchange', this.onHashchange);
		this.filter = window.location.hash.split("#/")[1];
	},
	beforeDestroy: function() {
		window.removeEventListener('hashchange', this.onHashchange);
	},
	methods: {
		// raise events
		onToggle: function(id) {
			this.$emit('toggle', id);
		},
		onUpdate: function(todo) {
			this.$emit('update', todo);
		},
		onRemove: function(id) {
			this.$emit('remove', id);
		},
		onDrop: function(drop) {
			this.$emit('drop', drop);
		},
		// hashchange
		onHashchange: function (todos) {
			this.filter = window.location.hash.split("#/")[1];
		},
		active: function (todos) {
			return this.todos.filter(todo => {
				return !todo.completed;
			});
		},
		completed: function (todos) {
			return this.todos.filter(todo => {
				return todo.completed;
			});
		}
	},
	computed: {
		todosFiltered: {
			get: function() {
				if (this.filter == 'completed') {
					return this.completed();
				} else if (this.filter == 'active') {
					return this.active();
				}
				this.filter = 'all';
				return this.todos;
			}
		}
	}
}
</script>

<style scoped>
	.todo-list {
		// margin: 6px;
		border: 1px solid #8d600d00; /*orange*/
	}

	.filters {
		display: flex;
		align-items: center;
		justify-content: center;
		border: 1px solid #8d600d00; /*orange*/
		margin: 6px 0px;
		// padding: 0px 0px 6px 6px;
		font-size: 0.8em;
	}
	.filters a {
		border: 1px solid #25466c; /*blue*/
		margin: 0px 3px;
		padding: 3px 6px;
		color: inherit;
		text-decoration: none;
		border-radius: 5px;
		font-family: sans-serif;
	}
	a.selected {
		border-color: #474747;
		color: #666;
	}
</style>
