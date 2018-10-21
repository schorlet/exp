<template>
	<form @submit.prevent="onSubmit">
		<input
			autofocus
			autocomplete="off"
			placeholder="What needs to be done ?"
			name="title"
			v-model="title"
		/>
		<input type="submit" value="+"/>
	</form>
</template>

<script>
module.exports = {
	name: 'TodoInput',
	data: function() {
		return {
			title: ''
		}
	},
	methods: {
		onSubmit: function() {
			this.$emit('create', this.title);
			this.title = '';
			this.$el.title.focus();
		}
	}
}
</script>

<style scoped>
	form {
		display: flex;
		align-items: center;
		margin: 6px;
		border: 1px solid #8d600d; /*orange*/
	}
	input ,input:active, input:focus {
		margin: 6px;
		padding: 6px;
		border: 1px solid #8d0d8d; /*magenta*/
		font-size: inherit;
		line-height: inherit;
	}

	input[type=submit] {
		font-family: monospace;
	}
	input[name=title] {
		width: 100%;
		flex 1 1 auto;
	}
</style>
