<template>
	<article
		class="todo-item"
		:class="{completed: todo.completed, editing: editing, dragover: dragover}"
		:draggable="!editing"
		@dragstart="onDragstart"
		@dragend="onDragend"
		@dragover="onDragover"
		@dragleave="onDragleave"
		@drop.prevent="onDrop"
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
					<span v-html="highlighted()"></span>
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
		},
		highlight: {
			type: String,
			default: ''
		}
	},
	data: function() {
		return {
			editing: false,
			dragging: false,
			dragover: false
		}
	},
	updated: function () {
		if (this.editing) {
			this.$el.querySelector('.editable-input').focus();
		}
	},
	methods: {
		// drag source
		onDragstart: function(e) {
			e.dataTransfer.setData('text/plain', this.todo.id);
			this.dragging = true;
		},
		onDragend: function() {
			this.dragging = false;
		},
		// drag target
		onDragover: function(e) {
			if (this.dragging) {
				e.dropEffect = 'none';
			} else {
				this.dragover = true;
				e.dropEffect = 'move';
				e.preventDefault();
			}
		},
		onDragleave: function() {
			this.dragover = false;
		},
		onDrop: function(e) {
			const data = e.dataTransfer.getData("text/plain");
			this.dragover = false;
			this.$emit('drop', {
				from: data,
				to: this.todo.id
			});
		},
		// highlight
		highlighted: function() {
			if (!this.highlight) {
				return this.todo.title;
            }
            try {
				return this.todo.title.replace(
					new RegExp(`(${this.highlight})`, 'ig'),
					'<span class="highlight">$1</span>');
			} catch {
				return this.todo.title;
			}
		},
		// toggle
		onToggle: function() {
			this.$emit('toggle', this.todo.id);
		},
		// edit
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
		// remove
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
	.dragover {
		border: 1px dashed #8d0d0d; /*red*/
	}

	/* input,label */
	input, label, span {
		margin: 6px 6px 0px 0px;
		padding: 6px;
		font-size: inherit;
		line-height: inherit;
	}
	span, input[type=button] {
		font-family: monospace;
	}

	.highlight {
		background-color: #8d8d0d; /*yellow*/
		color: #1f2023; /*background*/
		margin: 0px;
		padding: 0px;
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
		border: 1px solid #333; /*magenta*/
		outline: none;
		flex: 1 1 auto;
		min-width: 0;
	}
	.editing .editable-input {
		display: inline-block;
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
		border: 1px solid #333; /*magenta*/
		margin: 6px 0px 0px 0px;
	}

	/* .editable-button */
	.editable-button {
		display: none;
		vertical-align: baseline;
		text-align: center;
		border: 1px solid #333; /*magenta*/
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
