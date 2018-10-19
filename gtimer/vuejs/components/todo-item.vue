<template>
	<article
		class="todo-item"
		:class="{completed: todo.completed}">

		<input
			type="checkbox"
			:checked="todo.completed"
			@click="onToggle"
		/>

		<input
			type="text"
			:value="todo.title"
			@blur="onEditUpdate"
			@keyup.enter="onEditUpdate"
			@keyup.esc="onEditUndo"
		/>

		<input
			type="button"
			class="destroy"
			@click="onRemove"
			value="&Cross;"
		/>

	</article>
</template>

<script>
module.exports = {
	name: 'TodoItem',
	props: {
		todo: {
			type: Object,
			required: true
		}
	},
	methods: {
		onToggle: function() {
			this.$emit('toggle', this.todo.id);
		},
		onEditUpdate: function(event) {
			this.$emit('update', {
				id: this.todo.id,
				title: event.target.value
			});
		},
		onEditUndo: function() {
			document.execCommand('undo');
		},
		onRemove: function() {
			this.$emit('remove', this.todo.id);
		}
	}
}
</script>

<style scoped>
	.todo-item input {
		margin: 6px;
		padding: 6px;
		border: 1px solid #ccc;
		font-size: 24px;
		line-height: 1.4em;
	}
	.todo-item input[type=text] {
		width: 70%;
	}

	.todo-item .destroy {
		display: inline;
	}
	.todo-item:hover .destroy {
		color: orange;
		display: inline;
	}
</style>
