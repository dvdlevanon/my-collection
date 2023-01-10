import Item from "./Item"

function ItemsList({ items }) {
	return (
		<div className="items">
		{items.map((item) => {
			return <div key={item.id}>
			<Item item={item} />
			</div>
		})}
		</div>
	)
}

export default ItemsList