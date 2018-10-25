<template>
	<article
		class="todo-item"
		:class="{completed: todo.completed, editing: this.editing}"
	>
		<span
			class="toggle"
			title="toggle"
			@click="onToggle"
		>
			&check;
		</span>

		<div class="editable">
			<div class="editable-label">
				<label
					@dblclick="onEditStart"
				>
					<span>{{todo.title}}</span>
				</label>

				<input
					type="button"
					class="editable-button"
					value="&odot;"
					title="edit"
					@click="onEditStart"
				/>
			</div>

			<input
				type="text"
				class="editable-input"
				:value="todo.title"
				@blur="onEditUpdate"
				@keyup.enter="onEditUpdate"
				@keyup.esc="onEditUndo"
			/>
		</div>

		<input
			type="button"
			value="&cross;"
			class="destroy"
			title="delete"
			@click="onRemove"
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
			this.$el.querySelector('.editable-input').focus();
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
			this.editing = false;
			if (this.todo.title === event.target.value) {
				return;
			}
			this.$emit('update', {
				id: this.todo.id,
				title: event.target.value
			});
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
		// margin: 6px;
		border: 1px solid #8d0d0d00; /*red*/
		border-bottom: 1px solid #8d0d0d; /*red*/
		padding: 0px 0px 6px 6px;
	}

	/* input,label */
	input, label, span {
		margin: 6px 6px 0px 0px;
		padding: 6px;
		font-size: inherit;
		line-height: inherit;
	}
	input {
		border: 1px solid #8d0d8d; /*magenta*/
	}
	span, input[type=button] {
		font-family: monospace;
	}

	/* .editing */
	.editable {
		display: flex;
		align-items: center;
		flex: 1 1 auto;
		min-width: 0;
	}
	.editable-input {
		display: none;
	}
	.editing .editable-input {
		display: inline-block;
		outline: none;
		flex: 1 1 auto;
		min-width: 0;
	}
	.editing .editable-label {
		display: none;
	}

	/* .editable-label */
	.editable-label {
		display: flex;
		align-items: center;
		flex: 1 1 auto;
		min-width: 0;
	}
	.editable-label label {
		border: 1px solid #8d0d8d00; /*magenta*/
		flex: 1 1 auto;
		min-width: 0;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}
	.editable-label:hover label {
		border: 1px solid #8d0d8d; /*magenta*/
		margin: 6px 0px 0px 0px;
	}

	/* .editable-button */
	.editable-button {
		display: none;
		vertical-align: baseline;
		text-align: center;
		border: 1px solid #8d0d8d; /*magenta*/
	}
	.editable-label:hover .editable-button {
		display: inline-block;
		border-left: 0px;
	}


	/* .toggle, .completed */
	.toggle {
		border: 1px solid #8d600d; /*orange*/
		// border-radius: 50%;
		vertical-align: baseline;
		text-align: center;
		width: 0.5em;
		height: 0.5em;
		line-height: 0.5em;
		color: #0d9d0d00; /*green*/
		cursor: default;
	}
	// .toggle:hover {
		// color: inherit;
	// }
	.completed .toggle {
		color: #0d9d0d; /*green*/
	}
	.completed label {
		text-decoration: line-through;
	}

	/* .destroy */
	.destroy {
		// border: 1px solid #8d600d; /*orange*/
		border: 0px;
		vertical-align: baseline;
		text-align: center;
		cursor: pointer;
	}
	.todo-item:hover .destroy {
		color: #8d0d0d; /*red*/
	}
</style>
