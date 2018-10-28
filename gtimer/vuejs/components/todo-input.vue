<template>
	<div class="todo-input">
		<span
			class="toggle"
			:class="{toggled: this.toggled}"
			title="toggle"
			@click="onToggle"
		>
			&check;
		</span>

		<form @submit.prevent="onSubmit">
			<input
				autofocus
				autocomplete="off"
				placeholder="What needs to be done ?"
				name="title"
				:value="title"
				@input="onInput"
			/>
			<input type="submit"
				title="create" value="&crarr;"/>
		</form>
	</div>
</template>

<script>
module.exports = {
	name: 'TodoInput',
	model: {
		prop: 'highlight',
		event: 'highlight'
	},
	props: {
		highlight: {
			type: String,
			default: ''
		}
	},
	data: function() {
		return {
			title: '',
			toggled: false,
		}
	},
	methods: {
		onToggle: function() {
			this.toggled = !this.toggled;
			this.$emit('toggle-all', this.toggled);
		},
		onInput: function(event) {
			this.title = event.target.value;
			this.$emit('highlight', this.title);
		},
		onSubmit: function() {
			this.$emit('create', this.title);
			this.title = '';
			this.$emit('highlight', this.title);
			this.$el.querySelector('input[name=title]').focus();
		}
	}
}
</script>

<style scoped>
	.todo-input {
		display: flex;
		align-items: center;
		padding: 0px 0px 6px 6px;
		border: 2px solid #8d600d00; /*orange*/
		border-bottom: 1px solid #8d600d; /*orange*/
	}

	/* .toggle */
	.toggle {
		margin: 6px 6px 0px 0px;
		padding: 6px;
		font-family: monospace;
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
	.toggled.toggle {
		color: #0d9d0d; /*green*/
	}

	form {
		display: flex;
		border: 1px solid #8d600d00; /*orange*/
		// padding: 0px 0px 6px 6px;
		width: 100%;
		flex 1 1 auto;
	}
	input {
		margin: 6px 6px 0px 0px;
		padding: 6px;
		border: 1px solid #8d0d8d00; /*magenta*/
		font-size: inherit;
		line-height: inherit;
	}

	input[type=submit] {
		font-family: monospace;
		cursor: pointer;
	}
	input[type=submit]:hover {
		color: #8d600d; /*orange*/
	}

	input[name=title] {
		width: 100%;
		flex 1 1 auto;
	}
	input[name=title]:focus {
		outline: none;
	}
</style>
