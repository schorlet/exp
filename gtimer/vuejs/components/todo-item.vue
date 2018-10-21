<template>
	<article
		class="todo-item"
		:class="{completed: todo.completed, editing: this.editing}"
	>
		<div class="view">
			<input
				type="button"
				class="toggle"
				value="&sext;"
				@click="onToggle"
			/>

			<div class="view-label">
				<label
					@dblclick="onEditStart"
				>
					{{todo.title}}
				</label>

				<input
					type="button"
					class="destroy"
					value="&Cross;"
					@click="onRemove"
				/>
			</div>
		</div>

		<input
			type="text"
			class="edit"
			:value="todo.title"
			@blur="onEditUpdate"
			@keyup.enter="onEditUpdate"
			@keyup.esc="onEditUndo"
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
	data: function() {
		return {
			editing: false,
		}
	},
	updated: function () {
		if (this.editing) {
			this.$el.querySelector('.edit').focus();
		}
	},
	methods: {
		onToggle: function() {
			this.$emit('toggle', this.todo.id);
		},
		onEditStart: function(event) {
			this.editing = true;
		},
		onEditUpdate: function(event) {
			this.$emit('update', {
				id: this.todo.id,
				title: event.target.value
			});
			this.editing = false;
		},
		onEditUndo: function() {
			document.execCommand('undo');
			this.editing = false;
		},
		onRemove: function() {
			this.$emit('remove', this.todo.id);
			this.editing = false;
		}
	}
}
</script>

<style scoped>
	.todo-item {
		display: flex;
		align-items: center;
		margin: 6px;
		border: 1px solid #8d0d0d;
	}

	/* input,label */
	input, label {
		margin: 6px;
		padding: 6px;
		border: 1px solid #8d0d8d; /*magenta*/
		font-size: inherit;
		line-height: inherit;
	}
	input[type=text],label {
		width: 100%;
		flex 1 1 auto;
	}
	input[type=button] {
		font-family: monospace;
	}

	/* .editing */
	.edit {
		display: none;
	}
	.view {
		display: flex;
		width: 100%;
	}
	.editing .edit {
		display: block;
	}
	.editing .view {
		display: none;
	}

	/* .completed */
	.completed label {
		text-decoration: line-through;
	}
	.completed .toggle {
		color: #8d600d; /*orange*/
	}

	/* .destroy */
	.view-label {
		display: flex;
		width: 100%;
	}
	.destroy {
		display: none;
	}
	.view-label:hover .destroy {
		color: #8d600d; /*orange*/
		display: block;
	}
</style>
