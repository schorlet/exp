new Vue({
	el: '#todo',
	components: {
		TodoApp: httpVueLoader('./components/todo-app.vue')
	}
});
