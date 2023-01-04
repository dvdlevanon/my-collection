import GalleryItem from "./GalleryItem"

function Gallery({ items }) {
  return (
    <div className="gallery">
      {items.map((item) => {
        return <div key={item.id}>
          <GalleryItem item={item} />
        </div>
      })}
    </div>
  )
}

export default Gallery