package hu.bme.aut.gallery

import android.content.Intent
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import hu.bme.aut.gallery.databinding.ActivityFolderBinding

class FolderActivity : AppCompatActivity() {

    private lateinit var binding: ActivityFolderBinding

    override fun onCreate(savedInstanceState: Bundle?) {

        // hide action bar
        supportActionBar?.hide()

        super.onCreate(savedInstanceState)
        binding = ActivityFolderBinding.inflate(layoutInflater)
        setContentView(binding.root)


        binding.galleryButton.setOnClickListener {
            startActivity(Intent(this, GalleryActivity::class.java))
        }
        binding.calendarButton.setOnClickListener {
            startActivity(Intent(this, CalendarActivity::class.java))
        }

    }
}