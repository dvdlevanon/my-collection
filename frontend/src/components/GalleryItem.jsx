function GalleryItem({ item }) {
  return (
    <div className="gallery-item">
        <img src={item.cover} alt='' />
        <span className="item-title">{item.title}</span>
    </div>
  )
}

export default GalleryItem