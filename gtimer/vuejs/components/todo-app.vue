<template>
	<section class="todoapp" v-cloak>
		<header>
			<todo-input
				@create="onCreate"
				@toggle-all="onToggleAll"
				v-model="highlight"
			></todo-input>
		</header>

		<section>
			<todo-list
				:todos="todos"
				@toggle="onToggle"
				@update="onUpdate"
				@remove="onRemove"
				@drop="onDrop"
				:highlight="highlight"
			></todo-list>
		</section>

		<footer>
			{{todos}}
		</footer>
	</section>
</template>

<script>
module.exports = {
	name: 'TodoApp',
	components: {
		TodoInput: httpVueLoader('./todo-input.vue'),
		TodoList: httpVueLoader('./todo-list.vue')
	},
	data: function() {
		return {
			debug: true,
			todos: [
				{id: '1', title:'text 1', completed: true},
				{id: '2', title:''},
				{id: '3', title:'text 3 is longer and should be ellipsed'}
			],
			count: 3,
			highlight: ''
		}
	},
	methods: {
		log: function(message) {
			if (this.debug) console.log(message);
		},
		onCreate: function(title) {
			this.log(`onCreate: title:${title}`);
			this.count++;
			this.todos.push({
				id: this.count.toString(),
				title: title
			});
		},
		onToggle: function(id) {
			this.log(`onToggle: id:${id}`);
			const index = this.todos.findIndex(item => item.id === id);
			if (index >= 0) {
				const todo = this.todos[index];
				const completed = todo.completed || false;
				this.$set(todo, 'completed', !completed);
			}
		},
		onToggleAll: function(state) {
			this.log(`onToggleAll: state:${state}`);
			this.todos.forEach(todo => {
				this.$set(todo, 'completed', state);
			});
		},
		onUpdate: function(updated) {
			this.log(`onUpdate: id:${updated.id}, title:${updated.title}`);
			const index = this.todos.findIndex(item => item.id === updated.id);
			if (index >= 0) {
				const todo = this.todos[index];
				todo.title = updated.title;
			}
		},
		onRemove: function(id) {
			this.log(`onRemove: id:${id}`);
			const index = this.todos.findIndex(item => item.id === id);
			if (index >= 0) {
				this.todos.splice(index, 1);
			}
		},
		onDrop: function(drop) {
			this.log(`onDrop: from:${drop.from} to:${drop.to}`);
			const ifrom = this.todos.findIndex(item => item.id === drop.from);
			const ito = this.todos.findIndex(item => item.id === drop.to);
			if (ifrom >= 0 && ito >= 0) {
				const tfrom = this.todos.splice(ifrom, 1);
				this.todos.splice(ito, 0, tfrom[0]);
			}
		},
	}
}
</script>

<style scoped>
	[v-cloak] {
	  display: none;
	}

	.todoapp {
		margin: 30px 6px;
		border: 1px solid #0d9d0d; /*green*/
	}

	header,section,footer {
		margin: 0px 6px;
		border: 1px solid #8d8d0d00; /*yellow*/
	}
	footer {
		line-height: 1.1em;
		white-space: pre;
		font-family: monospace;
		font-size: 12px;
	}
</style>
