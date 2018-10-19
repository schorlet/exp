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
	input {
		padding: 6px;
		margin: 6px;
		border: 1px solid #ccc;
		font-size: 24px;
		line-height: 1.4em;
	}

	input[name=title] {
		width: 80%;
	}
</style>
