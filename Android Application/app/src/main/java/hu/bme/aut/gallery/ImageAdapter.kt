package hu.bme.aut.gallery

import android.R
import android.content.Context
import android.content.Intent
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import androidx.recyclerview.widget.RecyclerView
import com.squareup.picasso.Picasso
import java.io.File

class ImageAdapter(private val mList: List<ImageModel>): RecyclerView.Adapter<ImageAdapter.RecyclerViewHolder>() {


    private val context: Context? = null

    private val imagePathArrayList: ArrayList<String>? = null


    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int) : RecyclerViewHolder {


        val view = LayoutInflater.from(parent.context).inflate(R.layout.card_layout, parent, false)

        return RecyclerViewHolder(view)
    }
    override fun onBindViewHolder(holder: RecyclerViewHolder, position: Int) {
        val imgFile = File(imagePathArrayList?.get(position).toString() ?: String())

        if (imgFile.exists()) {

            // if the file exists then we are displaying that file in our image view using picasso library.
            Picasso.get().load(imgFile).placeholder(R.drawable.gallery_thumb)
                .into(holder.imageIV)

            // on below line we are adding click listener to our item of recycler view.
            holder.itemView.setOnClickListener { // inside on click listener we are creating a new intent
                val i = Intent(context, SingelPicActivity::class.java)

                // on below line we are passing the image path to our new activity.
                i.putExtra("imgPath", imagePathArrayList?.get(position))

                // at last we are starting our activity.
                context?.startActivity(i)
            }
        }
    }


    override fun getItemCount(): Int = imagePathArrayList?.size!!

    class RecyclerViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        // creating variables for our views.
        val imageIV: ImageView

        init {
            // initializing our views with their ids.
            imageIV = itemView.findViewById(R.id.idIVImage)
        }
    }

    inner class RecycleViewHolder(val binding: ImageBinding) : RecyclerView.ViewHolder(binding.root)

}
