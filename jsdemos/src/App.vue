<template>
	<div id="app">
		<div class="list">
			<button
				v-for="(item,index) in items"
				:key="index"
				class="item"
				@click="clickOntheTap(item)"
			>{{item.name}}</button>
		</div>
	</div>
</template>

<script>
import HelloWorld from "./components/HelloWorld.vue";
import { sxw } from "./ex/ex1";
import { suanfa } from "./ex/ex2";
export default {
	name: "app",
	data() {
		return {
			items: [{ unt: sxw, name: "例子1" }, { unt: suanfa, name: "算法" }]
		};
	},
	components: {
		HelloWorld
	},
	created() {
		// this.test();
	},
	methods: {
		clickOntheTap(item) {
			// console.log(item);
			item.unt();
		}
	}
};
</script>

<style lang="scss">
#app {
	font-family: "Avenir", Helvetica, Arial, sans-serif;
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
	text-align: center;
	color: #2c3e50;
	margin-top: 60px;
}
.list {
	width: 100%;
	display: -webkit-flex; /* Safari */
	display: flex;
	flex-direction: column;
	// justify-content: space-between;
	align-items: center;
	.item {
		margin-top: 20px;
		width: 80%;
		height: 60px;
		background-color: brown;
		color: white;
		font-size: 30px;
	}
}
</style>>